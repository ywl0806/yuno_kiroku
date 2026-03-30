package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/ywl0806/yuno_kiroku/internal/consts"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/enums"
	"github.com/ywl0806/yuno_kiroku/internal/store"
	imagepkg "github.com/ywl0806/yuno_kiroku/pkg/image"
)

// ResizeService MinIO webhook으로 수신한 원본 이미지를 리사이즈하고
// view/thumbnail을 생성한 뒤 face_recognition_job을 등록합니다.
type ResizeService struct {
	mediaItemStore store.MediaItemStore
	imageUploader  *ImageUploader
	faceDispatcher FaceRecognitionDispatcher
}

func NewResizeService(
	mediaItemStore store.MediaItemStore,
	imageUploader *ImageUploader,
	faceDispatcher FaceRecognitionDispatcher,
) *ResizeService {
	return &ResizeService{
		mediaItemStore: mediaItemStore,
		imageUploader:  imageUploader,
		faceDispatcher: faceDispatcher,
	}
}

// ProcessResize 원본 이미지 키를 받아 리사이즈 파이프라인 전체를 처리합니다.
func (s *ResizeService) ProcessResize(ctx context.Context, originalKey string) error {
	// 1. storage_key로 media_item_id 조회
	mediaItemID, err := s.mediaItemStore.GetMediaItemIDByStorageKey(ctx, originalKey)
	if err != nil {
		return fmt.Errorf("media_item_id 조회 실패 (key=%s): %w", originalKey, err)
	}

	// 2. 처리 중으로 상태 변경
	if _, err = s.mediaItemStore.UpdateMediaItemUploadStatus(ctx, db.UpdateMediaItemUploadStatusParams{
		ID:           mediaItemID,
		UploadStatus: string(enums.UploadStatusProcessing),
	}); err != nil {
		return fmt.Errorf("status 업데이트 실패: %w", err)
	}

	// 3. 원본 파일 다운로드
	originalData, err := s.imageUploader.storage.GetFile(originalKey)
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

	// 6. 리사이즈 (view, thumbnail)
	viewImg, err := imagepkg.Resize(originalData, consts.VIEW_MAX_LENGTH)
	if err != nil {
		s.setFailed(ctx, mediaItemID, err)
		return fmt.Errorf("view 리사이즈 실패: %w", err)
	}
	thumbImg, err := imagepkg.Resize(originalData, consts.THUMBNAIL_MAX_LENGTH)
	if err != nil {
		s.setFailed(ctx, mediaItemID, err)
		return fmt.Errorf("thumbnail 리사이즈 실패: %w", err)
	}
	originalData = nil

	// 7. 업로드 경로 결정
	// originalKey = "{familyId}/{albumId}/{date}/original/{uuid}.{ext}"
	uploadFolder := extractStorageKeyFolder(originalKey)
	baseName := extractStorageKeyBaseName(originalKey)

	// 8. view 업로드
	viewKey, err := s.imageUploader.storage.SaveFile(
		viewImg.Data, uploadFolder,
		consts.VIEW_STORAGE_PREFIX+"/"+baseName+viewImg.Ext,
	)
	viewImg.Data = nil
	if err != nil {
		s.setFailed(ctx, mediaItemID, err)
		return fmt.Errorf("view 업로드 실패: %w", err)
	}

	// 9. thumbnail 업로드
	thumbKey, err := s.imageUploader.storage.SaveFile(
		thumbImg.Data, uploadFolder,
		consts.THUMBNAIL_STORAGE_PREFIX+"/"+baseName+thumbImg.Ext,
	)
	thumbImg.Data = nil
	if err != nil {
		s.setFailed(ctx, mediaItemID, err)
		return fmt.Errorf("thumbnail 업로드 실패: %w", err)
	}

	// 10. media_files 레코드 생성 (view, thumbnail)
	if _, err = s.mediaItemStore.CreateMediaFile(ctx, db.CreateMediaFileParams{
		MediaItemID: mediaItemID,
		Role:        string(enums.MediaItemRoleView),
		StorageKey:  viewKey,
		Width:       sql.NullInt32{Int32: int32(viewImg.Width), Valid: true},
		Height:      sql.NullInt32{Int32: int32(viewImg.Height), Valid: true},
	}); err != nil {
		s.setFailed(ctx, mediaItemID, err)
		return fmt.Errorf("view media_file 생성 실패: %w", err)
	}
	if _, err = s.mediaItemStore.CreateMediaFile(ctx, db.CreateMediaFileParams{
		MediaItemID: mediaItemID,
		Role:        string(enums.MediaItemRoleThumbnail),
		StorageKey:  thumbKey,
		Width:       sql.NullInt32{Int32: int32(thumbImg.Width), Valid: true},
		Height:      sql.NullInt32{Int32: int32(thumbImg.Height), Valid: true},
	}); err != nil {
		s.setFailed(ctx, mediaItemID, err)
		return fmt.Errorf("thumbnail media_file 생성 실패: %w", err)
	}

	// 11. upload_status = completed (업로드 프로세스 완료)
	if _, err = s.mediaItemStore.UpdateMediaItemUploadStatus(ctx, db.UpdateMediaItemUploadStatusParams{
		ID:           mediaItemID,
		UploadStatus: string(enums.UploadStatusCompleted),
	}); err != nil {
		log.Printf("upload_status completed 업데이트 실패 (무시): %v", err)
	}

	// 12. face_recognition_job 생성 (얼굴 인식 프로세스 시작)
	mediaItem, err := s.mediaItemStore.GetMediaItemByID(ctx, mediaItemID)
	if err != nil {
		log.Printf("media_item 조회 실패, face_recognition_job 생략: %v", err)
	} else {
		if err = s.faceDispatcher.Dispatch(ctx, db.CreateFaceRecognitionJobParams{
			MediaItemID:    mediaItemID,
			FamilyID:       mediaItem.FamilyID,
			ViewStorageKey: viewKey,
		}); err != nil {
			log.Printf("face_recognition_job 디스패치 실패 (무시): %v", err)
		}
	}

	log.Printf("리사이즈 처리 완료: media_item_id=%d, view=%s, thumb=%s", mediaItemID, viewKey, thumbKey)
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

// "1/2/2025-03-27/original/uuid.jpg" → "1/2/2025-03-27"
func extractStorageKeyFolder(key string) string {
	parts := splitStorageKey(key)
	if len(parts) < 4 {
		return ""
	}
	return parts[0] + "/" + parts[1] + "/" + parts[2]
}

// "1/2/2025-03-27/original/uuid.jpg" → "uuid" (확장자 없는 파일명)
func extractStorageKeyBaseName(key string) string {
	parts := splitStorageKey(key)
	if len(parts) == 0 {
		return "file"
	}
	name := parts[len(parts)-1]
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			return name[:i]
		}
	}
	return name
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
