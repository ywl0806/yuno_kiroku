package services

import (
	"context"
	"database/sql"

	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store"
)

type TagService struct {
	tagStore store.TagStore
}

func NewTagService(tagStore store.TagStore) *TagService {
	return &TagService{tagStore: tagStore}
}

func (s *TagService) GetTagsByFamilyID(ctx context.Context, familyID int32) ([]db.Tag, error) {
	return s.tagStore.GetTagsByFamilyID(ctx, familyID)
}

func (s *TagService) CreateTag(ctx context.Context, familyID, userID int32, name string) (db.Tag, error) {
	return s.tagStore.CreateTag(ctx, db.CreateTagParams{
		FamilyID:  sql.NullInt32{Int32: familyID, Valid: true},
		Name:      name,
		CreatedBy: sql.NullInt32{Int32: userID, Valid: true},
	})
}

func (s *TagService) DeleteTag(ctx context.Context, tagID, familyID int32) error {
	return s.tagStore.DeleteTag(ctx, tagID, familyID)
}

func (s *TagService) AddTagToMediaItem(ctx context.Context, mediaItemID, tagID, userID int32) error {
	return s.tagStore.AddTagToMediaItem(ctx, db.AddTagToMediaItemParams{
		MediaItemID: mediaItemID,
		TagID:       tagID,
		TaggedBy:    userID,
	})
}

func (s *TagService) RemoveTagFromMediaItem(ctx context.Context, mediaItemID, tagID int32) error {
	return s.tagStore.RemoveTagFromMediaItem(ctx, db.RemoveTagFromMediaItemParams{
		MediaItemID: mediaItemID,
		TagID:       tagID,
	})
}

func (s *TagService) GetTagsForMediaItem(ctx context.Context, mediaItemID int32) ([]db.Tag, error) {
	return s.tagStore.GetTagsForMediaItem(ctx, mediaItemID)
}
