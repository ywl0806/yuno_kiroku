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
	GetMediaItemsByTakenAt(ctx context.Context, arg db.GetMediaItemsByTakenAtParams) ([]db.GetMediaItemsByTakenAtRow, error)
	GetMediaItemsByTakenAtWithIdentity(ctx context.Context, arg db.GetMediaItemsByTakenAtWithIdentityParams) ([]db.GetMediaItemsByTakenAtWithIdentityRow, error)
	GetMediaItemRange(ctx context.Context, clanGroupID int32) ([]db.GetMediaItemRangeRow, error)
	CreateUploadBatch(ctx context.Context, albumID int32) (db.UploadBatch, error)
	UpdateMediaItemUploadStatus(ctx context.Context, arg db.UpdateMediaItemUploadStatusParams) (db.UpdateMediaItemUploadStatusRow, error)
	GetUploadStatuses(ctx context.Context, uploadBatchID int32) ([]db.GetUploadStatusesRow, error)
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

func (s *mediaItemStore) GetMediaItemsByTakenAtWithIdentity(ctx context.Context, arg db.GetMediaItemsByTakenAtWithIdentityParams) ([]db.GetMediaItemsByTakenAtWithIdentityRow, error) {
	return wrapErr(s.queries.GetMediaItemsByTakenAtWithIdentity(ctx, arg))
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
