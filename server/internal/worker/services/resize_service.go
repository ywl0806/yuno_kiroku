package services

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"strings"

	"github.com/ywl0806/yuno_kiroku/internal/consts"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/enums"
	"github.com/ywl0806/yuno_kiroku/internal/services"
	"github.com/ywl0806/yuno_kiroku/internal/store"
	"github.com/ywl0806/yuno_kiroku/internal/utils"
	imagepkg "github.com/ywl0806/yuno_kiroku/pkg/image"
	"github.com/ywl0806/yuno_kiroku/pkg/storage"
)

// ResizeService MinIO webhook으로 수신한 원본 이미지를 리사이즈하고
// view/thumbnail을 생성한 뒤 face_recognition_job을 등록합니다.
// 비디오 파일의 경우 video-processing SQS 큐에 job을 발행합니다.
type ResizeService struct {
	log             *slog.Logger
	mediaItemStore  store.MediaItemStore
	storage         storage.StorageService
	faceDispatcher  services.FaceRecognitionDispatcher
	videoDispatcher services.VideoJobDispatcher
}

func NewResizeService(
	mediaItemStore store.MediaItemStore,
	storageService storage.StorageService,
	faceDispatcher services.FaceRecognitionDispatcher,
	videoDispatcher services.VideoJobDispatcher,
) *ResizeService {
	return &ResizeService{
		log:             slog.Default().With("layer", "worker", "component", "resize"),
		mediaItemStore:  mediaItemStore,
		storage:         storageService,
		faceDispatcher:  faceDispatcher,
		videoDispatcher: videoDispatcher,
	}
}

// ProcessResize 원본 파일 키를 받아 파이프라인 전체를 처리합니다.
// 비디오 파일이면 video-processing SQS job을 발행하고, 이미지이면 리사이즈 처리합니다.
func (s *ResizeService) ProcessResize(ctx context.Context, originalKey string) error {
	// 1. 키에서 mediaItemID / familyID 파싱 (key = "{familyId}/original/{mediaItemId}.ext")
	mediaItemID, err := extractMediaItemIDFromKey(originalKey)
	if err != nil {
		s.log.ErrorContext(ctx, "media_item_id 파싱 실패", "key", originalKey, "error", err)
		return err
	}
	familyID, err := extractFamilyIDFromKey(originalKey)
	if err != nil {
		s.log.ErrorContext(ctx, "family_id 파싱 실패", "key", originalKey, "error", err)
		return err
	}

	// S2-08: pending(01) → processing(02) 원자적 전환
	// 이미 processing 이상인 경우 (중복 이벤트) 즉시 반환
	ok, err := s.mediaItemStore.UpdateMediaItemToProcessingIfPending(ctx, mediaItemID)
	if err != nil {
		s.log.ErrorContext(ctx, "status 업데이트 실패", "media_item_id", mediaItemID, "error", err)
		return err
	}
	if !ok {
		s.log.InfoContext(ctx, "duplicate event skipped", "media_item_id", mediaItemID)
		return nil
	}

	// 비디오 파일이면 video-processing SQS job 발행 후 반환
	if isVideoKey(originalKey) {
		s.log.InfoContext(ctx, "video file detected, dispatching video job", "key", originalKey)
		ext := extractStorageKeyExt(originalKey)
		return s.videoDispatcher.Dispatch(ctx, services.VideoJobParams{
			MediaItemID:        mediaItemID,
			FamilyID:           familyID,
			OriginalStorageKey: originalKey,
			FileName:           mediaItemID + "." + ext,
			MimeType:           videoMimeType(ext),
		})
	}

	// 3. 파일 크기 사전 검증 — 다운로드 전 HeadObject로 확인하여 OOM 방지
	fileSize, err := s.storage.GetFileSize(ctx, originalKey)
	if err != nil {
		s.setFailedTransient(ctx, mediaItemID, err)
		return err
	}
	if fileSize > consts.MAX_IMAGE_FILE_SIZE {
		s.log.WarnContext(ctx, "file size too large", "media_item_id", mediaItemID, "file_size", fileSize)
		s.setFailed(ctx, mediaItemID, originalKey, enums.FailureReasonFileTooLarge, err)
		return nil
	}

	// 4. 원본 파일 다운로드
	originalData, err := s.storage.GetFile(ctx, originalKey)
	if err != nil {
		s.log.ErrorContext(ctx, "original file download failed", "media_item_id", mediaItemID, "error", err)
		s.setFailedTransient(ctx, mediaItemID, err)
		return err
	}

	// 5. 이미지 파싱 (EXIF 추출) — 영구 실패 시 nil 반환하여 SQS 재시도 방지
	ext := extractStorageKeyExt(originalKey)
	meta, err := imagepkg.Parse(originalData, ext)
	if err != nil {
		reason := classifyParseError(ext)
		s.log.ErrorContext(ctx, "image parse failed", "media_item_id", mediaItemID, "error", err)
		s.setFailed(ctx, mediaItemID, originalKey, reason, err)
		return nil
	}

	// 6. taken_at을 EXIF 값으로 업데이트
	// S2-06: EXIF 없는 파일(정상)과 DB 오류(이상)를 로그 레벨로 분리
	var nullLat, nullLon sql.NullFloat64
	if meta.Lat != nil {
		nullLat = sql.NullFloat64{Float64: *meta.Lat, Valid: true}
	}
	if meta.Lon != nil {
		nullLon = sql.NullFloat64{Float64: *meta.Lon, Valid: true}
	}
	if err = s.mediaItemStore.UpdateMediaItemTakenAt(ctx, db.UpdateMediaItemTakenAtParams{
		ID:                     mediaItemID,
		TakenAt:                meta.TakenAt,
		TakenLocationLatitude:  nullLat,
		TakenLocationLongitude: nullLon,
	}); err != nil {
		if meta.TakenAt.IsZero() {
			s.log.InfoContext(ctx, "no EXIF taken_at, keeping default", "media_item_id", mediaItemID)
		} else {
			s.log.ErrorContext(ctx, "taken_at DB update failed", "media_item_id", mediaItemID, "error", err)
		}
	}

	// 7. 리사이즈 → DB 저장 → 업로드 → completed → face job 디스패치
	if err = s.ProcessResizeFromData(ctx, originalData, mediaItemID, familyID); err != nil {
		s.setFailedTransient(ctx, mediaItemID, err)
		return err
	}

	s.log.InfoContext(ctx, "resize completed", "media_item_id", mediaItemID)
	return nil
}

// ProcessResizeFromData 이미 메모리에 있는 원본 데이터를 받아
// 리사이즈 → DB media_files 레코드 생성 → S3 업로드 → upload_status=completed → face job 디스패치를 처리합니다.
// multipart 업로드(MediaItemService)와 MinIO webhook(ProcessResize) 양쪽에서 공유합니다.
func (s *ResizeService) ProcessResizeFromData(ctx context.Context, originalData []byte, mediaItemID, familyID string) error {
	// 1. 리사이즈
	viewImg, err := imagepkg.Resize(originalData, consts.VIEW_MAX_LENGTH)
	if err != nil {
		s.log.ErrorContext(ctx, "view 리사이즈 실패", "media_item_id", mediaItemID, "error", err)
		return err
	}
	thumbImg, err := imagepkg.Resize(originalData, consts.THUMBNAIL_MAX_LENGTH)
	if err != nil {
		s.log.ErrorContext(ctx, "thumbnail 리사이즈 실패", "media_item_id", mediaItemID, "error", err)
		return err
	}
	originalData = nil

	// 2. storage key 사전 계산 (SaveFile은 입력 key를 그대로 반환)
	viewKey := utils.BuildMediaKey(familyID, consts.VIEW_STORAGE_PREFIX, mediaItemID, viewImg.Ext)
	thumbKey := utils.BuildMediaKey(familyID, consts.THUMBNAIL_STORAGE_PREFIX, mediaItemID, thumbImg.Ext)

	// 3. DB media_files 레코드 생성
	if _, err = s.mediaItemStore.CreateMediaFile(ctx, db.CreateMediaFileParams{
		MediaItemID: mediaItemID,
		Role:        string(enums.MediaItemRoleView),
		StorageKey:  viewKey,
		Width:       sql.NullInt32{Int32: int32(viewImg.Width), Valid: true},
		Height:      sql.NullInt32{Int32: int32(viewImg.Height), Valid: true},
	}); err != nil {
		s.log.ErrorContext(ctx, "view media_file DB 생성 실패", "media_item_id", mediaItemID, "error", err)
		return err
	}
	if _, err = s.mediaItemStore.CreateMediaFile(ctx, db.CreateMediaFileParams{
		MediaItemID: mediaItemID,
		Role:        string(enums.MediaItemRoleThumbnail),
		StorageKey:  thumbKey,
		Width:       sql.NullInt32{Int32: int32(thumbImg.Width), Valid: true},
		Height:      sql.NullInt32{Int32: int32(thumbImg.Height), Valid: true},
	}); err != nil {
		// view DB 레코드 롤백
		if delErr := s.mediaItemStore.DeleteMediaFileByItemAndRole(ctx, mediaItemID, string(enums.MediaItemRoleView)); delErr != nil {
			s.log.ErrorContext(ctx, "view media_file DB rollback failed", "media_item_id", mediaItemID, "error", delErr)
		}
		s.log.ErrorContext(ctx, "thumbnail media_file DB 생성 실패", "media_item_id", mediaItemID, "error", err)
		return err
	}

	// 4. view S3 업로드
	//    실패 시 DB 레코드 롤백 (S2-05)
	if _, err = s.storage.SaveFile(ctx, viewKey, viewImg.Data); err != nil {
		s.rollbackMediaFiles(ctx, mediaItemID)
		s.log.ErrorContext(ctx, "view S3 업로드 실패", "media_item_id", mediaItemID, "error", err)
		return err
	}
	viewImg.Data = nil

	// 5. thumbnail S3 업로드
	//    실패 시 view S3 롤백 + DB 레코드 롤백 (S2-04 + S2-05)
	if _, err = s.storage.SaveFile(ctx, thumbKey, thumbImg.Data); err != nil {
		if delErr := s.storage.DeleteFile(ctx, viewKey); delErr != nil {
			s.log.ErrorContext(ctx, "view S3 rollback failed — orphan file", "key", viewKey, "error", delErr)
		}
		s.rollbackMediaFiles(ctx, mediaItemID)
		s.log.ErrorContext(ctx, "thumbnail S3 업로드 실패", "media_item_id", mediaItemID, "error", err)
		return err
	}
	thumbImg.Data = nil

	// 6. upload_status = completed
	if _, err = s.mediaItemStore.UpdateMediaItemUploadStatus(ctx, db.UpdateMediaItemUploadStatusParams{
		ID:           mediaItemID,
		UploadStatus: string(enums.UploadStatusCompleted),
	}); err != nil {
		s.log.WarnContext(ctx, "upload_status completed update failed (ignored)", "media_item_id", mediaItemID, "error", err)
	}

	// 7. face_recognition_job 디스패치 (S2-07)
	//    실패 시 face_recognition_status = pending 유지 (배치 재발행 대상으로 추적됨)
	if err = s.faceDispatcher.Dispatch(ctx, services.FaceRecognitionJobParams{
		MediaItemID:    mediaItemID,
		FamilyID:       familyID,
		ViewStorageKey: viewKey,
	}); err != nil {
		s.log.ErrorContext(ctx, "face recognition job dispatch failed", "media_item_id", mediaItemID, "error", err)
	} else {
		if statusErr := s.mediaItemStore.UpdateFaceRecognitionStatus(ctx, mediaItemID, string(enums.FaceRecognitionStatusDispatched)); statusErr != nil {
			s.log.WarnContext(ctx, "face_recognition_status dispatched update failed", "media_item_id", mediaItemID, "error", statusErr)
		}
	}

	return nil
}

// rollbackMediaFiles view/thumbnail DB media_files 레코드를 삭제한다.
// S3 업로드 실패 후 DB 상태를 원상복구하기 위해 사용한다.
func (s *ResizeService) rollbackMediaFiles(ctx context.Context, mediaItemID string) {
	for _, role := range []string{string(enums.MediaItemRoleView), string(enums.MediaItemRoleThumbnail)} {
		if err := s.mediaItemStore.DeleteMediaFileByItemAndRole(ctx, mediaItemID, role); err != nil {
			s.log.ErrorContext(ctx, "media_file DB rollback failed", "media_item_id", mediaItemID, "role", role, "error", err)
		}
	}
}

// setFailed 영구 실패 처리 — failure_reason 기록 + S3 원본 파일 즉시 삭제
// 호출 후 nil을 반환해야 SQS 메시지가 delete되어 재시도가 발생하지 않음
func (s *ResizeService) setFailed(ctx context.Context, mediaItemID, originalKey string, reason enums.FailureReason, err error) {
	if updateErr := s.mediaItemStore.UpdateMediaItemFailed(ctx, mediaItemID, string(reason)); updateErr != nil {
		s.log.ErrorContext(ctx, "status failed update failed", "media_item_id", mediaItemID, "error", updateErr)
	}
	if delErr := s.storage.DeleteFile(ctx, originalKey); delErr != nil {
		s.log.ErrorContext(ctx, "original file delete failed (ignored)", "key", originalKey, "error", delErr)
	}
	s.log.ErrorContext(ctx, "permanent failure", "media_item_id", mediaItemID, "reason", reason, "error", err)
}

// setFailedTransient 일시적 실패 처리 — status만 '04'로 변경, SQS 재시도 허용
func (s *ResizeService) setFailedTransient(ctx context.Context, mediaItemID string, err error) {
	if _, updateErr := s.mediaItemStore.UpdateMediaItemUploadStatus(ctx, db.UpdateMediaItemUploadStatusParams{
		ID:           mediaItemID,
		UploadStatus: string(enums.UploadStatusFailed),
	}); updateErr != nil {
		s.log.ErrorContext(ctx, "status failed update failed", "media_item_id", mediaItemID, "error", updateErr)
	}
	s.log.WarnContext(ctx, "transient failure", "media_item_id", mediaItemID, "error", err)
}

var supportedImageExts = map[string]bool{
	"jpg": true, "jpeg": true, "png": true,
	"webp": true, "heic": true, "heif": true,
	"gif": true, "tiff": true, "tif": true,
}

func classifyParseError(ext string) enums.FailureReason {
	if !supportedImageExts[strings.ToLower(ext)] {
		return enums.FailureReasonUnsupportedFormat
	}
	return enums.FailureReasonCorruptedFile
}

// "1/2/2025-03-27/original/uuid.jpg" → "jpg"
func extractStorageKeyExt(key string) string {
	for i := len(key) - 1; i >= 0; i-- {
		if key[i] == '.' {
			return key[i+1:]
		}
		if key[i] == '/' {
			break
		}
	}
	return "jpg"
}

// "familyId/original/456.jpg" → familyId
func extractFamilyIDFromKey(key string) (string, error) {
	parts := splitStorageKey(key)
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid key format: %s", key)
	}
	return parts[0], nil
}

// "original/1/456.jpg" → mediaItemId="456"
func extractMediaItemIDFromKey(key string) (string, error) {
	parts := splitStorageKey(key)
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid key format: %s", key)
	}
	name := parts[len(parts)-1]
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			return name[:i], nil
		}
	}
	return "", fmt.Errorf("no extension in key: %s", key)
}

var videoExtensions = map[string]string{
	"mp4":  "video/mp4",
	"mov":  "video/quicktime",
	"avi":  "video/x-msvideo",
	"mkv":  "video/x-matroska",
	"webm": "video/webm",
	"m4v":  "video/x-m4v",
}

func isVideoKey(key string) bool {
	ext := extractStorageKeyExt(key)
	_, ok := videoExtensions[ext]
	return ok
}

func videoMimeType(ext string) string {
	if mime, ok := videoExtensions[ext]; ok {
		return mime
	}
	return "video/mp4"
}

func splitStorageKey(key string) []string {
	var parts []string
	start := 0
	for i, c := range key {
		if c == '/' {
			if i > start {
				parts = append(parts, key[start:i])
			}
			start = i + 1
		}
	}
	if start < len(key) {
		parts = append(parts, key[start:])
	}
	return parts
}
