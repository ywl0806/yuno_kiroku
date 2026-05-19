package services

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/ywl0806/yuno_kiroku/internal/apperr"
	"github.com/ywl0806/yuno_kiroku/internal/consts"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/enums"
	"github.com/ywl0806/yuno_kiroku/internal/store"
	internalutils "github.com/ywl0806/yuno_kiroku/internal/utils"
	"github.com/ywl0806/yuno_kiroku/pkg/storage"
)

type MediaItemService struct {
	mediaItemStore            store.MediaItemStore
	albumGroupPermissionStore store.AlbumGroupPermissionStore
	storage                   storage.StorageService
}

func NewMediaItemService(
	mediaItemStore store.MediaItemStore,
	albumGroupPermissionStore store.AlbumGroupPermissionStore,
	storageService storage.StorageService,
) *MediaItemService {
	return &MediaItemService{
		mediaItemStore:            mediaItemStore,
		albumGroupPermissionStore: albumGroupPermissionStore,
		storage:                   storageService,
	}
}

// UploadImageResult는 이미지 업로드 시작 후 반환하는 응답
type UploadImageResult struct {
	MediaItemID string
	Status      string
}

// 미디어 아이템의 업로드 상태를 업데이트
func (s *MediaItemService) UpdateMediaItemUploadStatus(ctx context.Context, mediaItemID string, uploadStatus enums.UploadStatus) error {
	_, err := s.mediaItemStore.UpdateMediaItemUploadStatus(ctx, db.UpdateMediaItemUploadStatusParams{
		ID:           mediaItemID,
		UploadStatus: string(uploadStatus),
	})
	return err
}

// GetMediaItemsByTakenAt는 촬영 시간 범위 내의 미디어 아이템을 반환
func (s *MediaItemService) GetMediaItemsByTakenAt(ctx context.Context, groupID int32, userID string, from, to time.Time) ([]db.GetMediaItemsByTakenAtRow, error) {
	params := db.GetMediaItemsByTakenAtParams{
		GroupID:     groupID,
		TakenAtFrom: from,
		TakenAtTo:   to,
		UserID:      userID,
	}
	mediaItems, err := s.mediaItemStore.GetMediaItemsByTakenAt(ctx, params)
	if err != nil {
		log.Println("get media items by taken at error:", err)
		return nil, err
	}
	return mediaItems, nil
}

// MediaItemRange는 미디어 아이템이 존재하는 연/월 조합
type MediaItemRange struct {
	Year  int `json:"year"`
	Month int `json:"month"`
}

// GetMediaItemRange는 미디어 아이템이 있는 연/월 목록을 반환
func (s *MediaItemService) GetMediaItemRange(ctx context.Context, groupId int32) ([]MediaItemRange, error) {
	ranges, err := s.mediaItemStore.GetMediaItemRange(ctx, groupId)
	if err != nil {
		return nil, err
	}

	mediaItemRanges := make([]MediaItemRange, len(ranges))
	for i, r := range ranges {
		year, err := strconv.Atoi(r.Year)
		if err != nil {
			return nil, err
		}
		month, err := strconv.Atoi(r.Month)
		if err != nil {
			return nil, err
		}
		mediaItemRanges[i] = MediaItemRange{Year: year, Month: month}
	}

	return mediaItemRanges, nil
}

const SearchPageSize = 30

type SearchMediaItemsResult struct {
	Items   []db.SearchMediaItemsRow
	HasNext bool
}

// SearchMediaItems는 검색 조건으로 미디어 아이템을 페이지 단위로 반환
func (s *MediaItemService) SearchMediaItems(ctx context.Context, groupID int32, userID string, from, to *time.Time, albumID *string, identityIDs []int32, liked bool, tagIDs []int32, page int) (*SearchMediaItemsResult, error) {
	if identityIDs == nil {
		identityIDs = []int32{}
	}
	if tagIDs == nil {
		tagIDs = []int32{}
	}
	albumIDNull := uuid.NullUUID{Valid: false}
	if albumID != nil {
		parsed, err := uuid.Parse(*albumID)
		if err != nil {
			return nil, err
		}
		albumIDNull = uuid.NullUUID{UUID: parsed, Valid: true}
	}
	fromNull := sql.NullTime{Valid: false}
	if from != nil {
		fromNull = sql.NullTime{Time: *from, Valid: true}
	}
	toNull := sql.NullTime{Valid: false}
	if to != nil {
		toNull = sql.NullTime{Time: *to, Valid: true}
	}
	if page < 1 {
		page = 1
	}
	offset := int32((page - 1) * SearchPageSize)
	// 1件多く取得してhas_nextを判定
	params := db.SearchMediaItemsParams{
		GroupID:     groupID,
		UserID:      userID,
		TakenAtFrom: fromNull,
		TakenAtTo:   toNull,
		AlbumID:     albumIDNull,
		IdentityIds: identityIDs,
		Liked:       liked,
		TagIds:      tagIDs,
		PageOffset:  offset,
		PageSize:    int32(SearchPageSize) + 1,
	}
	items, err := s.mediaItemStore.SearchMediaItems(ctx, params)
	if err != nil {
		return nil, err
	}
	hasNext := len(items) > SearchPageSize
	if hasNext {
		items = items[:SearchPageSize]
	}
	return &SearchMediaItemsResult{Items: items, HasNext: hasNext}, nil
}

// PresignedUploadResult는 Presigned URL 발급 후 반환하는 응답
type PresignedUploadResult struct {
	MediaItemID   string
	PresignedURL  string
	StorageKey    string
	OriginalIndex int
}

var allowedContentTypes = map[string]bool{
	"image/jpeg": true, "image/png": true, "image/webp": true,
	"image/heic": true, "image/heif": true, "image/gif": true,
	"image/tiff": true, "image/bmp": true,
	"video/mp4": true, "video/quicktime": true, "video/x-msvideo": true,
	"video/x-matroska": true, "video/webm": true, "video/x-m4v": true,
}

// CreatePresignedUpload는 S3 직접 업로드용 Presigned PUT URL을 발급
// Presigned URL 발급 실패 시 생성된 DB 레코드를 수동으로 롤백함
func (s *MediaItemService) CreatePresignedUpload(
	ctx context.Context,
	fileName, contentType string,
	familyId, albumId string, uploadBatchID int32,
) (*PresignedUploadResult, error) {
	if !allowedContentTypes[contentType] {
		return nil, apperr.NewValidationError("message.validation.unsupported_format", nil)
	}

	// 1. media_item 생성 (storageKey에 mediaItemID 필요)
	mediaItem, err := s.mediaItemStore.CreateMediaItem(ctx, db.CreateMediaItemParams{
		FamilyID:      familyId,
		AlbumID:       albumId,
		UploadBatchID: uploadBatchID,
		TakenAt:       time.Now(), // Resize Worker가 EXIF 파싱 후 실제 값으로 업데이트
		FileName:      sql.NullString{String: fileName, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	// 2. mediaItemID 확정 후 키 생성 ({familyId}/original/{mediaItemId}.ext)
	storageKey := internalutils.BuildMediaKeyFromFileName(familyId, consts.ORIGINAL_STORAGE_PREFIX, mediaItem.ID, fileName)

	_, err = s.mediaItemStore.CreateMediaFile(ctx, db.CreateMediaFileParams{
		MediaItemID: mediaItem.ID,
		Role:        string(enums.MediaItemRoleOriginal),
		StorageKey:  storageKey,
	})
	if err != nil {
		if rollbackErr := s.mediaItemStore.DeleteMediaItem(ctx, mediaItem.ID); rollbackErr != nil {
			log.Printf("media_item 롤백 실패 (id=%s): %v", mediaItem.ID, rollbackErr)
		}
		return nil, err
	}

	// 3. Presigned URL 발급 — 실패 시 DB 레코드 삭제 (media_files는 CASCADE)
	presignedURL, err := s.storage.GeneratePresignedPutURL(ctx, storageKey, contentType, time.Hour)
	if err != nil {
		if rollbackErr := s.mediaItemStore.DeleteMediaItem(ctx, mediaItem.ID); rollbackErr != nil {
			log.Printf("media_item 롤백 실패 (id=%s): %v", mediaItem.ID, rollbackErr)
		}
		return nil, err
	}

	return &PresignedUploadResult{
		MediaItemID:  mediaItem.ID,
		PresignedURL: presignedURL,
		StorageKey:   storageKey,
	}, nil
}

// BatchPresignedUploadItem는 배치 Presigned URL 발급 요청 항목
type BatchPresignedUploadItem struct {
	FileName    string
	ContentType string
}

// BatchPresignedUploadResult는 배치 Presigned URL 발급 결과 (부분 성공 지원)
type BatchPresignedUploadResult struct {
	Success []*PresignedUploadResult
	Failed  []BatchPresignedUploadFailedItem
}

// BatchPresignedUploadFailedItem은 배치 발급 실패 항목
type BatchPresignedUploadFailedItem struct {
	FileName string
	Index    int
	Reason   string // "unsupported_format" 등
}

// CreateBatchPresignedUpload는 여러 파일에 대한 Presigned PUT URL을 한 번에 발급
// 개별 항목 실패 시 해당 항목 DB 레코드는 트랜잭션 롤백으로 자동 정리되며, 나머지 항목은 계속 처리
func (s *MediaItemService) CreateBatchPresignedUpload(
	ctx context.Context,
	items []BatchPresignedUploadItem,
	familyId, albumId string,
	uploadBatchID int32,
) (*BatchPresignedUploadResult, error) {
	result := &BatchPresignedUploadResult{
		Success: make([]*PresignedUploadResult, 0, len(items)),
		Failed:  make([]BatchPresignedUploadFailedItem, 0),
	}
	for i, item := range items {
		r, err := s.CreatePresignedUpload(ctx, item.FileName, item.ContentType, familyId, albumId, uploadBatchID)
		if err != nil {
			log.Printf("batch presigned upload 실패 (index=%d, file=%s): %v", i, item.FileName, err)
			reason := ""
			if !allowedContentTypes[item.ContentType] {
				reason = string(enums.FailureReasonUnsupportedFormat)
			}
			result.Failed = append(result.Failed, BatchPresignedUploadFailedItem{FileName: item.FileName, Index: i, Reason: reason})
			continue
		}
		r.OriginalIndex = i
		result.Success = append(result.Success, r)
	}
	return result, nil
}

// CreateUploadBatch는 앨범에 대한 새 업로드 배치를 생성
func (s *MediaItemService) CreateUploadBatch(ctx context.Context, familyId, albumId string) (*db.UploadBatch, error) {
	uploadBatch, err := s.mediaItemStore.CreateUploadBatch(ctx, albumId)
	if err != nil {
		return nil, err
	}
	return &uploadBatch, nil
}

// GetUploadStatuses는 배치 내 모든 아이템의 업로드 상태를 반환합니다.
func (s *MediaItemService) GetUploadStatuses(ctx context.Context, uploadBatchID int32) ([]db.GetUploadStatusesRow, error) {
	return s.mediaItemStore.GetUploadStatuses(ctx, uploadBatchID)
}

type BatchThumbnail struct {
	ID              string
	ThumbnailUrl    string
	ThumbnailWidth  int32
	ThumbnailHeight int32
}

type UploadBatchWithThumbnails struct {
	ID         int32
	AlbumID    string
	UploadAt   time.Time
	Count      int64
	Thumbnails []BatchThumbnail
}

// GetUploadBatchesWithThumbnails는 배치 목록과 각 배치의 썸네일 5개를 반환합니다. (2-쿼리 방식)
func (s *MediaItemService) GetUploadBatchesWithThumbnails(ctx context.Context, groupID int32, page int) ([]UploadBatchWithThumbnails, bool, error) {
	if page < 1 {
		page = 1
	}
	offset := int32((page - 1) * SearchPageSize)

	// Step 1: 배치 목록 + 카운트 (페이지네이션, 완료된 아이템만 카운트)
	batches, err := s.mediaItemStore.GetUploadBatchesAndMediaItemCounts(ctx, db.GetUploadBatchesAndMediaItemCountsParams{
		GroupID:    groupID,
		PageSize:   int32(SearchPageSize) + 1,
		PageOffset: offset,
	})
	if err != nil {
		return nil, false, err
	}

	hasNext := len(batches) > SearchPageSize
	if hasNext {
		batches = batches[:SearchPageSize]
	}

	if len(batches) == 0 {
		return []UploadBatchWithThumbnails{}, hasNext, nil
	}

	// Step 2: 배치 ID 배열로 썸네일 한 번에 조회
	batchIDs := make([]int32, len(batches))
	for i, b := range batches {
		batchIDs[i] = b.ID
	}

	thumbnailRows, err := s.mediaItemStore.GetUploadBatchWithThumbnails(ctx, batchIDs)
	if err != nil {
		return nil, false, err
	}

	// Step 3: batch ID로 그룹화
	thumbnailMap := make(map[int32][]BatchThumbnail)
	for _, t := range thumbnailRows {
		width := int32(0)
		height := int32(0)
		if t.Width.Valid {
			width = t.Width.Int32
		}
		if t.Height.Valid {
			height = t.Height.Int32
		}
		thumbnailMap[t.ID] = append(thumbnailMap[t.ID], BatchThumbnail{
			ID:              t.MediaItemID,
			ThumbnailUrl:    internalutils.ParseStoragePath(t.StorageKey),
			ThumbnailWidth:  width,
			ThumbnailHeight: height,
		})
	}

	result := make([]UploadBatchWithThumbnails, len(batches))
	for i, b := range batches {
		result[i] = UploadBatchWithThumbnails{
			ID:         b.ID,
			AlbumID:    b.AlbumID,
			UploadAt:   b.UploadAt,
			Count:      b.Count,
			Thumbnails: thumbnailMap[b.ID],
		}
	}
	return result, hasNext, nil
}

type UploadBatchItemsResult struct {
	Items   []db.GetMediaItemsByUploadBatchIdRow
	HasNext bool
}

// GetMediaItemsByUploadBatch는 특정 배치의 미디어 아이템을 페이지 단위로 반환합니다.
func (s *MediaItemService) GetMediaItemsByUploadBatch(ctx context.Context, uploadBatchID int32, userID string, page int) (*UploadBatchItemsResult, error) {
	if page < 1 {
		page = 1
	}
	offset := int32((page - 1) * SearchPageSize)
	items, err := s.mediaItemStore.GetMediaItemsByUploadBatchId(ctx, db.GetMediaItemsByUploadBatchIdParams{
		UserID:        userID,
		UploadBatchID: uploadBatchID,
		PageOffset:    offset,
		PageSize:      int32(SearchPageSize) + 1,
	})
	if err != nil {
		return nil, err
	}
	hasNext := len(items) > SearchPageSize
	if hasNext {
		items = items[:SearchPageSize]
	}
	return &UploadBatchItemsResult{Items: items, HasNext: hasNext}, nil
}

// DeleteMediaItem는 미디어 아이템을 삭제합니다.
func (s *MediaItemService) DeleteMediaItem(ctx context.Context, mediaItemID string, familyID string, userID string) error {
	// media_item 조회
	mediaItem, err := s.mediaItemStore.GetMediaItemByIDAndFamilyID(ctx, mediaItemID, familyID)
	if err != nil {
		return err
	}

	hasPermission, err := s.albumGroupPermissionStore.CheckUserHasPermissionForAlbum(ctx, mediaItem.AlbumID, userID, "W")
	if err != nil {
		return err
	}

	if !hasPermission {
		return apperr.NewForbiddenError("error.forbidden", nil)
	}

	return s.mediaItemStore.DeleteMediaItem(ctx, mediaItemID)
}

// UpdateMediaItemAlbum는 미디어 아이템의 앨범을 업데이트합니다.
func (s *MediaItemService) UpdateMediaItemAlbum(ctx context.Context, id string, familyID string, userID string, albumID string) error {
	mediaItem, err := s.mediaItemStore.GetMediaItemByIDAndFamilyID(ctx, id, familyID)
	if err != nil {
		return err
	}

	hasPermission, err := s.albumGroupPermissionStore.CheckUserHasPermissionForAlbum(ctx, mediaItem.AlbumID, userID, "W")
	if err != nil {
		return err
	}

	if !hasPermission {
		return apperr.NewForbiddenError("error.forbidden", nil)
	}

	return s.mediaItemStore.UpdateMediaItemAlbum(ctx, db.UpdateMediaItemAlbumParams{
		ID:      id,
		AlbumID: albumID,
	})
}
