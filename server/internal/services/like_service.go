package services

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store"
)

type LikeService interface {
	LikeMediaItem(ctx context.Context, mediaItemID, userID int32) error
	UnlikeMediaItem(ctx context.Context, mediaItemID, userID int32) error
	IsMediaItemLiked(ctx context.Context, mediaItemID, userID int32) (bool, error)
}

type likeService struct {
	likeStore store.LikeStore
}

func NewLikeService(likeStore store.LikeStore) LikeService {
	return &likeService{likeStore: likeStore}
}

func (s *likeService) LikeMediaItem(ctx context.Context, mediaItemID, userID int32) error {
	return s.likeStore.LikeMediaItem(ctx, db.LikeMediaItemParams{
		MediaItemID: mediaItemID,
		UserID:      userID,
	})
}

func (s *likeService) UnlikeMediaItem(ctx context.Context, mediaItemID, userID int32) error {
	return s.likeStore.UnlikeMediaItem(ctx, db.UnlikeMediaItemParams{
		MediaItemID: mediaItemID,
		UserID:      userID,
	})
}

func (s *likeService) IsMediaItemLiked(ctx context.Context, mediaItemID, userID int32) (bool, error) {
	return s.likeStore.IsMediaItemLiked(ctx, db.IsMediaItemLikedParams{
		MediaItemID: mediaItemID,
		UserID:      userID,
	})
}
