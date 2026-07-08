package services

import (
	"context"
	"database/sql"

	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store"
)

type TagService interface {
	GetTagsByFamilyID(ctx context.Context, familyID int32) ([]db.Tag, error)
	CreateTag(ctx context.Context, familyID, userID int32, name string) (db.Tag, error)
	DeleteTag(ctx context.Context, tagID, familyID int32) error
	AddTagToMediaItem(ctx context.Context, mediaItemID, tagID, userID int32) error
	RemoveTagFromMediaItem(ctx context.Context, mediaItemID, tagID int32) error
	GetTagsForMediaItem(ctx context.Context, mediaItemID int32) ([]db.Tag, error)
}

type tagService struct {
	tagStore store.TagStore
}

func NewTagService(tagStore store.TagStore) TagService {
	return &tagService{tagStore: tagStore}
}

func (s *tagService) GetTagsByFamilyID(ctx context.Context, familyID int32) ([]db.Tag, error) {
	return s.tagStore.GetTagsByFamilyID(ctx, familyID)
}

func (s *tagService) CreateTag(ctx context.Context, familyID, userID int32, name string) (db.Tag, error) {
	return s.tagStore.CreateTag(ctx, db.CreateTagParams{
		FamilyID:  sql.NullInt32{Int32: familyID, Valid: true},
		Name:      name,
		CreatedBy: sql.NullInt32{Int32: userID, Valid: true},
	})
}

func (s *tagService) DeleteTag(ctx context.Context, tagID, familyID int32) error {
	return s.tagStore.DeleteTag(ctx, tagID, familyID)
}

func (s *tagService) AddTagToMediaItem(ctx context.Context, mediaItemID, tagID, userID int32) error {
	return s.tagStore.AddTagToMediaItem(ctx, db.AddTagToMediaItemParams{
		MediaItemID: mediaItemID,
		TagID:       tagID,
		TaggedBy:    userID,
	})
}

func (s *tagService) RemoveTagFromMediaItem(ctx context.Context, mediaItemID, tagID int32) error {
	return s.tagStore.RemoveTagFromMediaItem(ctx, db.RemoveTagFromMediaItemParams{
		MediaItemID: mediaItemID,
		TagID:       tagID,
	})
}

func (s *tagService) GetTagsForMediaItem(ctx context.Context, mediaItemID int32) ([]db.Tag, error) {
	return s.tagStore.GetTagsForMediaItem(ctx, mediaItemID)
}
