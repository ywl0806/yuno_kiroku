package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/ywl0806/yuno_kiroku/internal/db"
)

// TagStore 태그 데이터 접근 인터페이스
type TagStore interface {
	GetTagsByFamilyID(ctx context.Context, familyID string) ([]db.Tag, error)
	CreateTag(ctx context.Context, arg db.CreateTagParams) (db.Tag, error)
	DeleteTag(ctx context.Context, id int32, familyID string) error
	AddTagToMediaItem(ctx context.Context, arg db.AddTagToMediaItemParams) error
	RemoveTagFromMediaItem(ctx context.Context, arg db.RemoveTagFromMediaItemParams) error
	GetTagsForMediaItem(ctx context.Context, mediaItemID string) ([]db.Tag, error)
}

type tagStore struct {
	queries *db.Queries
}

func NewTagStore(queries *db.Queries) TagStore {
	return &tagStore{queries: queries}
}

func (s *tagStore) GetTagsByFamilyID(ctx context.Context, familyID string) ([]db.Tag, error) {
	uid, err := uuid.Parse(familyID)
	if err != nil {
		return nil, err
	}
	return wrapErr(s.queries.GetTagsByFamilyID(ctx, uuid.NullUUID{UUID: uid, Valid: true}))
}

func (s *tagStore) CreateTag(ctx context.Context, arg db.CreateTagParams) (db.Tag, error) {
	return wrapErr(s.queries.CreateTag(ctx, arg))
}

func (s *tagStore) DeleteTag(ctx context.Context, id int32, familyID string) error {
	uid, err := uuid.Parse(familyID)
	if err != nil {
		return err
	}
	return s.queries.DeleteTag(ctx, db.DeleteTagParams{
		ID:       id,
		FamilyID: uuid.NullUUID{UUID: uid, Valid: true},
	})
}

func (s *tagStore) AddTagToMediaItem(ctx context.Context, arg db.AddTagToMediaItemParams) error {
	return s.queries.AddTagToMediaItem(ctx, arg)
}

func (s *tagStore) RemoveTagFromMediaItem(ctx context.Context, arg db.RemoveTagFromMediaItemParams) error {
	return s.queries.RemoveTagFromMediaItem(ctx, arg)
}

func (s *tagStore) GetTagsForMediaItem(ctx context.Context, mediaItemID string) ([]db.Tag, error) {
	return wrapErr(s.queries.GetTagsForMediaItem(ctx, mediaItemID))
}
