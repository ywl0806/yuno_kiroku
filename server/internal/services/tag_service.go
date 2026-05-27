package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store"
)

type TagService struct {
	tagStore store.TagStore
}

func NewTagService(tagStore store.TagStore) *TagService {
	return &TagService{tagStore: tagStore}
}

func (s *TagService) GetTagsByFamilyID(ctx context.Context, familyID string) ([]db.Tag, error) {
	return s.tagStore.GetTagsByFamilyID(ctx, familyID)
}

func (s *TagService) CreateTag(ctx context.Context, familyID, userID string, name string) (db.Tag, error) {
	familyUUID, err := uuid.Parse(familyID)
	if err != nil {
		return db.Tag{}, err
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return db.Tag{}, err
	}
	return s.tagStore.CreateTag(ctx, db.CreateTagParams{
		FamilyID:  uuid.NullUUID{UUID: familyUUID, Valid: true},
		Name:      name,
		CreatedBy: uuid.NullUUID{UUID: userUUID, Valid: true},
	})
}

func (s *TagService) DeleteTag(ctx context.Context, tagID int32, familyID string) error {
	return s.tagStore.DeleteTag(ctx, tagID, familyID)
}

func (s *TagService) AddTagToMediaItem(ctx context.Context, mediaItemID string, tagID int32, userID string) error {
	return s.tagStore.AddTagToMediaItem(ctx, db.AddTagToMediaItemParams{
		MediaItemID: mediaItemID,
		TagID:       tagID,
		TaggedBy:    userID,
	})
}

func (s *TagService) RemoveTagFromMediaItem(ctx context.Context, mediaItemID string, tagID int32) error {
	return s.tagStore.RemoveTagFromMediaItem(ctx, db.RemoveTagFromMediaItemParams{
		MediaItemID: mediaItemID,
		TagID:       tagID,
	})
}

func (s *TagService) GetTagsForMediaItem(ctx context.Context, mediaItemID string) ([]db.Tag, error) {
	return s.tagStore.GetTagsForMediaItem(ctx, mediaItemID)
}
