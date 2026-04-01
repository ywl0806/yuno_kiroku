package services

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/enums"
	"github.com/ywl0806/yuno_kiroku/internal/store"
)

type MediaItemService struct {
	mediaItemStore store.MediaItemStore
	imageUploader  *ImageUploader
}

func NewMediaItemService(
	mediaItemStore store.MediaItemStore,
	imageUploader *ImageUploader,
) *MediaItemService {
	return &MediaItemService{
		mediaItemStore: mediaItemStore,
		imageUploader:  imageUploader,
	}
}

// UploadImageResult는 이미지 업로드 시작 후 반환하는 응답
type UploadImageResult struct {
	MediaItemID int32
	Status      string
}

// 미디어 아이템의 업로드 상태를 업데이트
func (s *MediaItemService) UpdateMediaItemUploadStatus(ctx context.Context, mediaItemID int32, uploadStatus enums.UploadStatus) error {
	_, err := s.mediaItemStore.UpdateMediaItemUploadStatus(ctx, db.UpdateMediaItemUploadStatusParams{
		ID:           mediaItemID,
		UploadStatus: string(uploadStatus),
	})
	return err
}

// GetMediaItemsByTakenAtHome는 필터 없이 촬영 시간 범위 내의 미디어 아이템을 반환 (홈 전용 경량 쿼리)
func (s *MediaItemService) GetMediaItemsByTakenAtHome(ctx context.Context, groupID int32, from, to time.Time) ([]db.GetMediaItemsByTakenAtHomeRow, error) {
	params := db.GetMediaItemsByTakenAtHomeParams{
		GroupID:     groupID,
		TakenAtFrom: from,
		TakenAtTo:   to,
	}
	mediaItems, err := s.mediaItemStore.GetMediaItemsByTakenAtHome(ctx, params)
	if err != nil {
		log.Println("get media items by taken at home error:", err)
		return nil, err
	}
	return mediaItems, nil
}

// GetMediaItemsByTakenAt는 촬영 시간 범위 내의 미디어 아이템을 반환
func (s *MediaItemService) GetMediaItemsByTakenAt(ctx context.Context, groupID int32, from, to time.Time, albumID *int32, identityIDs []int32) ([]db.GetMediaItemsByTakenAtRow, error) {
	if identityIDs == nil {
		identityIDs = []int32{}
	}
	albumIDNull := sql.NullInt32{Valid: false}
	if albumID != nil {
		albumIDNull = sql.NullInt32{Int32: *albumID, Valid: true}
	}
	params := db.GetMediaItemsByTakenAtParams{
		GroupID:     groupID,
		TakenAtFrom: from,
		TakenAtTo:   to,
		AlbumID:     albumIDNull,
		IdentityIds: identityIDs,
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
func (s *MediaItemService) SearchMediaItems(ctx context.Context, groupID int32, from, to *time.Time, albumID *int32, identityIDs []int32, page int) (*SearchMediaItemsResult, error) {
	if identityIDs == nil {
		identityIDs = []int32{}
	}
	albumIDNull := sql.NullInt32{Valid: false}
	if albumID != nil {
		albumIDNull = sql.NullInt32{Int32: *albumID, Valid: true}
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
		TakenAtFrom: fromNull,
		TakenAtTo:   toNull,
		AlbumID:     albumIDNull,
		IdentityIds: identityIDs,
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
	MediaItemID  int32
	PresignedURL string
	StorageKey   string
}

// CreatePresignedUpload는 S3 직접 업로드용 Presigned PUT URL을 발급
// DB 레코드를 먼저 생성(DB-first)하여 고아 파일을 방지
func (s *MediaItemService) CreatePresignedUpload(
	ctx context.Context,
	fileName, contentType string,
	familyId, albumId, uploadBatchID int32,
) (*PresignedUploadResult, error) {
	// 1. media_item 먼저 생성 (storageKey에 mediaItemID 필요)
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

	// 2. mediaItemID 확정 후 키 생성 (original/{familyId}/{mediaItemId}.ext)
	storageKey := s.imageUploader.BuildOriginalKey(familyId, mediaItem.ID, fileName)

	_, err = s.mediaItemStore.CreateMediaFile(ctx, db.CreateMediaFileParams{
		MediaItemID: mediaItem.ID,
		Role:        string(enums.MediaItemRoleOriginal),
		StorageKey:  storageKey,
	})
	if err != nil {
		return nil, err
	}

	presignedURL, err := s.imageUploader.GeneratePresignedPutURL(ctx, storageKey, contentType, time.Hour)
	if err != nil {
		return nil, err
	}

	return &PresignedUploadResult{
		MediaItemID:  mediaItem.ID,
		PresignedURL: presignedURL,
		StorageKey:   storageKey,
	}, nil
}

// CreateUploadBatch는 앨범에 대한 새 업로드 배치를 생성
func (s *MediaItemService) CreateUploadBatch(ctx context.Context, familyId, albumId int32) (*db.UploadBatch, error) {
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
