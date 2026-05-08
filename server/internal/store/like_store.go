package store

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

// LikeStore 좋아요 데이터 접근 인터페이스
type LikeStore interface {
	LikeMediaItem(ctx context.Context, arg db.LikeMediaItemParams) error
	UnlikeMediaItem(ctx context.Context, arg db.UnlikeMediaItemParams) error
	IsMediaItemLiked(ctx context.Context, arg db.IsMediaItemLikedParams) (bool, error)
}

type likeStore struct {
	queries *db.Queries
}

func NewLikeStore(queries *db.Queries) LikeStore {
	return &likeStore{queries: queries}
}

func (s *likeStore) LikeMediaItem(ctx context.Context, arg db.LikeMediaItemParams) error {
	return s.queries.LikeMediaItem(ctx, arg)
}

func (s *likeStore) UnlikeMediaItem(ctx context.Context, arg db.UnlikeMediaItemParams) error {
	return s.queries.UnlikeMediaItem(ctx, arg)
}

func (s *likeStore) IsMediaItemLiked(ctx context.Context, arg db.IsMediaItemLikedParams) (bool, error) {
	return s.queries.IsMediaItemLiked(ctx, arg)
}
