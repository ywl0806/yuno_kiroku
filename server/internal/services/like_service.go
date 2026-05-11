package services

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store"
)

type LikeService struct {
	likeStore store.LikeStore
}

func NewLikeService(likeStore store.LikeStore) *LikeService {
	return &LikeService{likeStore: likeStore}
}

func (s *LikeService) LikeMediaItem(ctx context.Context, mediaItemID, userID string) error {
	return s.likeStore.LikeMediaItem(ctx, db.LikeMediaItemParams{
		MediaItemID: mediaItemID,
		UserID:      userID,
	})
}

func (s *LikeService) UnlikeMediaItem(ctx context.Context, mediaItemID, userID string) error {
	return s.likeStore.UnlikeMediaItem(ctx, db.UnlikeMediaItemParams{
		MediaItemID: mediaItemID,
		UserID:      userID,
	})
}

func (s *LikeService) IsMediaItemLiked(ctx context.Context, mediaItemID, userID string) (bool, error) {
	return s.likeStore.IsMediaItemLiked(ctx, db.IsMediaItemLikedParams{
		MediaItemID: mediaItemID,
		UserID:      userID,
	})
}
