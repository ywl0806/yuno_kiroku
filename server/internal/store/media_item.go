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
	GetMediaItemByIDAndFamilyID(ctx context.Context, id string, familyID string) (db.MediaItem, error)
	UpdateMediaItemTakenAt(ctx context.Context, arg db.UpdateMediaItemTakenAtParams) error
	GetMediaItemsByTakenAt(ctx context.Context, arg db.GetMediaItemsByTakenAtParams) ([]db.GetMediaItemsByTakenAtRow, error)
	GetMediaItemRange(ctx context.Context, clanGroupID int32) ([]db.GetMediaItemRangeRow, error)
	CreateUploadBatch(ctx context.Context, albumID string) (db.UploadBatch, error)
	UpdateMediaItemUploadStatus(ctx context.Context, arg db.UpdateMediaItemUploadStatusParams) (db.UpdateMediaItemUploadStatusRow, error)
	GetUploadStatuses(ctx context.Context, uploadBatchID int32) ([]db.GetUploadStatusesRow, error)
	SearchMediaItems(ctx context.Context, arg db.SearchMediaItemsParams) ([]db.SearchMediaItemsRow, error)
	GetUploadBatchesAndMediaItemCounts(ctx context.Context, arg db.GetUploadBatchesAndMediaItemCountsParams) ([]db.GetUploadBatchesAndMediaItemCountsRow, error)
	GetMediaItemsByUploadBatchId(ctx context.Context, arg db.GetMediaItemsByUploadBatchIdParams) ([]db.GetMediaItemsByUploadBatchIdRow, error)
	GetUploadBatchWithThumbnails(ctx context.Context, batchIds []int32) ([]db.GetUploadBatchWithThumbnailsRow, error)
	DeleteMediaItem(ctx context.Context, id string) error
	UpdateMediaItemAlbum(ctx context.Context, arg db.UpdateMediaItemAlbumParams) error
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

func (s *mediaItemStore) GetMediaItemRange(ctx context.Context, clanGroupID int32) ([]db.GetMediaItemRangeRow, error) {
	return wrapErr(s.queries.GetMediaItemRange(ctx, clanGroupID))
}

func (s *mediaItemStore) CreateUploadBatch(ctx context.Context, albumID string) (db.UploadBatch, error) {
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

func (s *mediaItemStore) GetMediaItemByIDAndFamilyID(ctx context.Context, id string, familyID string) (db.MediaItem, error) {
	return wrapErr(s.queries.GetMediaItemByIDAndFamilyID(ctx, db.GetMediaItemByIDAndFamilyIDParams{
		ID:       id,
		FamilyID: familyID,
	}))
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

func (s *mediaItemStore) DeleteMediaItem(ctx context.Context, id string) error {
	return mapDBError(s.queries.DeleteMediaItem(ctx, id))
}

func (s *mediaItemStore) UpdateMediaItemAlbum(ctx context.Context, arg db.UpdateMediaItemAlbumParams) error {
	return mapDBError(s.queries.UpdateMediaItemAlbum(ctx, arg))
}
