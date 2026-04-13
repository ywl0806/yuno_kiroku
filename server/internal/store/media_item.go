package store

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

// MediaItemStore 미디어 아이템/파일/업로드배치 데이터 접근 인터페이스
type MediaItemStore interface {
	CreateMediaItem(ctx context.Context, arg db.CreateMediaItemParams) (db.MediaItem, error)
	CreateMediaFile(ctx context.Context, arg db.CreateMediaFileParams) (db.MediaFile, error)
	GetMediaItemByFaceDetection(ctx context.Context, arg db.GetMediaItemByFaceDetectionParams) (db.MediaItem, error)
	GetMediaItemByID(ctx context.Context, id int32) (db.MediaItem, error)
	GetMediaItemIDByStorageKey(ctx context.Context, storageKey string) (int32, error)
	UpdateMediaItemTakenAt(ctx context.Context, arg db.UpdateMediaItemTakenAtParams) error
	GetMediaItemsByTakenAt(ctx context.Context, arg db.GetMediaItemsByTakenAtParams) ([]db.GetMediaItemsByTakenAtRow, error)
	GetMediaItemsByTakenAtHome(ctx context.Context, arg db.GetMediaItemsByTakenAtHomeParams) ([]db.GetMediaItemsByTakenAtHomeRow, error)
	GetMediaItemRange(ctx context.Context, clanGroupID int32) ([]db.GetMediaItemRangeRow, error)
	CreateUploadBatch(ctx context.Context, albumID int32) (db.UploadBatch, error)
	UpdateMediaItemUploadStatus(ctx context.Context, arg db.UpdateMediaItemUploadStatusParams) (db.UpdateMediaItemUploadStatusRow, error)
	GetUploadStatuses(ctx context.Context, uploadBatchID int32) ([]db.GetUploadStatusesRow, error)
	SearchMediaItems(ctx context.Context, arg db.SearchMediaItemsParams) ([]db.SearchMediaItemsRow, error)
	GetUploadBatchesAndMediaItemCounts(ctx context.Context, arg db.GetUploadBatchesAndMediaItemCountsParams) ([]db.GetUploadBatchesAndMediaItemCountsRow, error)
	GetMediaItemsByUploadBatchId(ctx context.Context, arg db.GetMediaItemsByUploadBatchIdParams) ([]db.GetMediaItemsByUploadBatchIdRow, error)
	GetUploadBatchWithThumbnails(ctx context.Context, batchIds []int32) ([]db.GetUploadBatchWithThumbnailsRow, error)
}

type mediaItemStore struct {
	queries *db.Queries
}

// NewMediaItemStore MediaItemStore 구현체 생성
func NewMediaItemStore(queries *db.Queries) MediaItemStore {
	return &mediaItemStore{queries: queries}
}

func (s *mediaItemStore) CreateMediaItem(ctx context.Context, arg db.CreateMediaItemParams) (db.MediaItem, error) {
	return wrapErr(s.queries.CreateMediaItem(ctx, arg))
}

func (s *mediaItemStore) CreateMediaFile(ctx context.Context, arg db.CreateMediaFileParams) (db.MediaFile, error) {
	return wrapErr(s.queries.CreateMediaFile(ctx, arg))
}

func (s *mediaItemStore) GetMediaItemByFaceDetection(ctx context.Context, arg db.GetMediaItemByFaceDetectionParams) (db.MediaItem, error) {
	return wrapErr(s.queries.GetMediaItemByFaceDetection(ctx, arg))
}

func (s *mediaItemStore) GetMediaItemsByTakenAt(ctx context.Context, arg db.GetMediaItemsByTakenAtParams) ([]db.GetMediaItemsByTakenAtRow, error) {
	return wrapErr(s.queries.GetMediaItemsByTakenAt(ctx, arg))
}

func (s *mediaItemStore) GetMediaItemsByTakenAtHome(ctx context.Context, arg db.GetMediaItemsByTakenAtHomeParams) ([]db.GetMediaItemsByTakenAtHomeRow, error) {
	return wrapErr(s.queries.GetMediaItemsByTakenAtHome(ctx, arg))
}

func (s *mediaItemStore) GetMediaItemRange(ctx context.Context, clanGroupID int32) ([]db.GetMediaItemRangeRow, error) {
	return wrapErr(s.queries.GetMediaItemRange(ctx, clanGroupID))
}

func (s *mediaItemStore) CreateUploadBatch(ctx context.Context, albumID int32) (db.UploadBatch, error) {
	return wrapErr(s.queries.CreateUploadBatch(ctx, albumID))
}

func (s *mediaItemStore) UpdateMediaItemUploadStatus(ctx context.Context, arg db.UpdateMediaItemUploadStatusParams) (db.UpdateMediaItemUploadStatusRow, error) {
	return wrapErr(s.queries.UpdateMediaItemUploadStatus(ctx, arg))
}

func (s *mediaItemStore) GetUploadStatuses(ctx context.Context, uploadBatchID int32) ([]db.GetUploadStatusesRow, error) {
	return wrapErr(s.queries.GetUploadStatuses(ctx, uploadBatchID))
}

func (s *mediaItemStore) SearchMediaItems(ctx context.Context, arg db.SearchMediaItemsParams) ([]db.SearchMediaItemsRow, error) {
	return wrapErr(s.queries.SearchMediaItems(ctx, arg))
}

func (s *mediaItemStore) GetMediaItemByID(ctx context.Context, id int32) (db.MediaItem, error) {
	return wrapErr(s.queries.GetMediaItemByID(ctx, id))
}

func (s *mediaItemStore) GetMediaItemIDByStorageKey(ctx context.Context, storageKey string) (int32, error) {
	return s.queries.GetMediaItemIDByStorageKey(ctx, storageKey)
}

func (s *mediaItemStore) UpdateMediaItemTakenAt(ctx context.Context, arg db.UpdateMediaItemTakenAtParams) error {
	return s.queries.UpdateMediaItemTakenAt(ctx, arg)
}

func (s *mediaItemStore) GetUploadBatchesAndMediaItemCounts(ctx context.Context, arg db.GetUploadBatchesAndMediaItemCountsParams) ([]db.GetUploadBatchesAndMediaItemCountsRow, error) {
	return wrapErr(s.queries.GetUploadBatchesAndMediaItemCounts(ctx, arg))
}

func (s *mediaItemStore) GetMediaItemsByUploadBatchId(ctx context.Context, arg db.GetMediaItemsByUploadBatchIdParams) ([]db.GetMediaItemsByUploadBatchIdRow, error) {
	return wrapErr(s.queries.GetMediaItemsByUploadBatchId(ctx, arg))
}

func (s *mediaItemStore) GetUploadBatchWithThumbnails(ctx context.Context, batchIds []int32) ([]db.GetUploadBatchWithThumbnailsRow, error) {
	return wrapErr(s.queries.GetUploadBatchWithThumbnails(ctx, batchIds))
}
