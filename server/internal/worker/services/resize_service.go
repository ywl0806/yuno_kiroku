package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/ywl0806/yuno_kiroku/internal/consts"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/enums"
	"github.com/ywl0806/yuno_kiroku/internal/services"
	"github.com/ywl0806/yuno_kiroku/internal/store"
	imagepkg "github.com/ywl0806/yuno_kiroku/pkg/image"
)

// ResizeService MinIO webhook으로 수신한 원본 이미지를 리사이즈하고
// view/thumbnail을 생성한 뒤 face_recognition_job을 등록합니다.
type ResizeService struct {
	mediaItemStore store.MediaItemStore
	imageUploader  *services.ImageUploader
	faceDispatcher services.FaceRecognitionDispatcher
}

func NewResizeService(
	mediaItemStore store.MediaItemStore,
	imageUploader *services.ImageUploader,
	faceDispatcher services.FaceRecognitionDispatcher,
) *ResizeService {
	return &ResizeService{
		mediaItemStore: mediaItemStore,
		imageUploader:  imageUploader,
		faceDispatcher: faceDispatcher,
	}
}

// ProcessResize 원본 이미지 키를 받아 리사이즈 파이프라인 전체를 처리합니다.
func (s *ResizeService) ProcessResize(ctx context.Context, originalKey string) error {
	// 1. 키에서 mediaItemID / familyID 파싱 (key = "original/{familyId}/{mediaItemId}.ext")
	mediaItemID, err := extractMediaItemIDFromKey(originalKey)
	if err != nil {
		return fmt.Errorf("media_item_id 파싱 실패 (key=%s): %w", originalKey, err)
	}
	familyID, err := extractFamilyIDFromKey(originalKey)
	if err != nil {
		return fmt.Errorf("family_id 파싱 실패 (key=%s): %w", originalKey, err)
	}

	// 2. 처리 중으로 상태 변경
	if _, err = s.mediaItemStore.UpdateMediaItemUploadStatus(ctx, db.UpdateMediaItemUploadStatusParams{
		ID:           mediaItemID,
		UploadStatus: string(enums.UploadStatusProcessing),
	}); err != nil {
		return fmt.Errorf("status 업데이트 실패: %w", err)
	}

	// 3. 원본 파일 다운로드
	originalData, err := s.imageUploader.GetFile(ctx, originalKey)
	if err != nil {
		s.setFailed(ctx, mediaItemID, err)
		return fmt.Errorf("원본 파일 다운로드 실패: %w", err)
	}

	// 4. 이미지 파싱 (EXIF 추출)
	ext := extractStorageKeyExt(originalKey)
	meta, err := imagepkg.Parse(originalData, ext)
	if err != nil {
		s.setFailed(ctx, mediaItemID, err)
		return fmt.Errorf("이미지 파싱 실패: %w", err)
	}

	// 5. taken_at을 EXIF 값으로 업데이트
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
		log.Printf("taken_at 업데이트 실패 (무시): %v", err)
	}

	// 6. 리사이즈 → 업로드 → media_files 저장 → completed → face job 디스패치
	if err = s.ProcessResizeFromData(ctx, originalData, mediaItemID, familyID); err != nil {
		s.setFailed(ctx, mediaItemID, err)
		return err
	}

	log.Printf("리사이즈 처리 완료: media_item_id=%d", mediaItemID)
	return nil
}

// ProcessResizeFromData 이미 메모리에 있는 원본 데이터를 받아
// 리사이즈 → 업로드 → media_files 저장 → upload_status=completed → face job 디스패치를 처리합니다.
// multipart 업로드(MediaItemService)와 MinIO webhook(ProcessResize) 양쪽에서 공유합니다.
func (s *ResizeService) ProcessResizeFromData(ctx context.Context, originalData []byte, mediaItemID, familyID int32) error {
	mediaItemStr := strconv.Itoa(int(mediaItemID))
	familyStr := strconv.Itoa(int(familyID))

	// 1. 리사이즈
	viewImg, err := imagepkg.Resize(originalData, consts.VIEW_MAX_LENGTH)
	if err != nil {
		return fmt.Errorf("view 리사이즈 실패: %w", err)
	}
	thumbImg, err := imagepkg.Resize(originalData, consts.THUMBNAIL_MAX_LENGTH)
	if err != nil {
		return fmt.Errorf("thumbnail 리사이즈 실패: %w", err)
	}
	originalData = nil

	// 2. view 업로드 (view/{familyId}/{mediaItemId}.ext)
	viewKey, err := s.imageUploader.SaveFile(
		ctx,
		viewImg.Data,
		consts.VIEW_STORAGE_PREFIX+"/"+familyStr,
		mediaItemStr+viewImg.Ext,
	)
	viewImg.Data = nil
	if err != nil {
		return fmt.Errorf("view 업로드 실패: %w", err)
	}

	// 3. thumbnail 업로드 (thumbnail/{familyId}/{mediaItemId}.ext)
	thumbKey, err := s.imageUploader.SaveFile(
		ctx,
		thumbImg.Data,
		consts.THUMBNAIL_STORAGE_PREFIX+"/"+familyStr,
		mediaItemStr+thumbImg.Ext,
	)
	thumbImg.Data = nil
	if err != nil {
		return fmt.Errorf("thumbnail 업로드 실패: %w", err)
	}

	// 4. media_files 레코드 생성
	if _, err = s.mediaItemStore.CreateMediaFile(ctx, db.CreateMediaFileParams{
		MediaItemID: mediaItemID,
		Role:        string(enums.MediaItemRoleView),
		StorageKey:  viewKey,
		Width:       sql.NullInt32{Int32: int32(viewImg.Width), Valid: true},
		Height:      sql.NullInt32{Int32: int32(viewImg.Height), Valid: true},
	}); err != nil {
		return fmt.Errorf("view media_file 생성 실패: %w", err)
	}
	if _, err = s.mediaItemStore.CreateMediaFile(ctx, db.CreateMediaFileParams{
		MediaItemID: mediaItemID,
		Role:        string(enums.MediaItemRoleThumbnail),
		StorageKey:  thumbKey,
		Width:       sql.NullInt32{Int32: int32(thumbImg.Width), Valid: true},
		Height:      sql.NullInt32{Int32: int32(thumbImg.Height), Valid: true},
	}); err != nil {
		return fmt.Errorf("thumbnail media_file 생성 실패: %w", err)
	}

	// 5. upload_status = completed
	if _, err = s.mediaItemStore.UpdateMediaItemUploadStatus(ctx, db.UpdateMediaItemUploadStatusParams{
		ID:           mediaItemID,
		UploadStatus: string(enums.UploadStatusCompleted),
	}); err != nil {
		log.Printf("upload_status completed 업데이트 실패 (무시): %v", err)
	}

	// 6. face_recognition_job 디스패치
	if err = s.faceDispatcher.Dispatch(ctx, services.FaceRecognitionJobParams{
		MediaItemID:    mediaItemID,
		FamilyID:       familyID,
		ViewStorageKey: viewKey,
	}); err != nil {
		log.Printf("face_recognition_job 디스패치 실패 (무시): %v", err)
	}

	return nil
}

func (s *ResizeService) setFailed(ctx context.Context, mediaItemID int32, err error) {
	if _, updateErr := s.mediaItemStore.UpdateMediaItemUploadStatus(ctx, db.UpdateMediaItemUploadStatusParams{
		ID:           mediaItemID,
		UploadStatus: string(enums.UploadStatusFailed),
	}); updateErr != nil {
		log.Printf("status failed 업데이트 실패: %v", updateErr)
	}
	log.Printf("media_item_id=%d 처리 실패: %v", mediaItemID, err)
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

// "original/1/456.jpg" → familyId=1
func extractFamilyIDFromKey(key string) (int32, error) {
	parts := splitStorageKey(key)
	if len(parts) < 3 {
		return 0, fmt.Errorf("invalid key format: %s", key)
	}
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("family_id 파싱 실패 (key=%s): %w", key, err)
	}
	return int32(id), nil
}

// "original/1/456.jpg" → mediaItemId=456
func extractMediaItemIDFromKey(key string) (int32, error) {
	parts := splitStorageKey(key)
	if len(parts) < 3 {
		return 0, fmt.Errorf("invalid key format: %s", key)
	}
	name := parts[len(parts)-1]
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			id, err := strconv.Atoi(name[:i])
			if err != nil {
				return 0, fmt.Errorf("media_item_id 파싱 실패 (key=%s): %w", key, err)
			}
			return int32(id), nil
		}
	}
	return 0, fmt.Errorf("no extension in key: %s", key)
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
