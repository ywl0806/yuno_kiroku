package store

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

// AlbumStore 앨범 데이터 접근 인터페이스
type AlbumStore interface {
	FindAlbumsForWrite(ctx context.Context, arg db.FindAlbumsForWriteParams) ([]db.Album, error)
}

type albumStore struct {
	queries *db.Queries
}

// NewAlbumStore AlbumStore 구현체 생성
func NewAlbumStore(queries *db.Queries) AlbumStore {
	return &albumStore{queries: queries}
}

func (s *albumStore) FindAlbumsForWrite(ctx context.Context, arg db.FindAlbumsForWriteParams) ([]db.Album, error) {
	return wrapErr(s.queries.FindAlbumsForWrite(ctx, arg))
}
