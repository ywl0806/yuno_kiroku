package store

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

// AlbumStore 앨범 데이터 접근 인터페이스
type AlbumStore interface {
	FindAlbumsForWrite(ctx context.Context, arg db.FindAlbumsForWriteParams) ([]db.Album, error)
	GetAlbumsOptions(ctx context.Context, familyID int32, groupID int32) ([]db.GetAlbumsOptionsRow, error)
	FindAlbumsByFamilyID(ctx context.Context, familyID int32) ([]db.Album, error)
	FindAlbumByIDAndFamilyID(ctx context.Context, id int32, familyID int32) (db.Album, error)
	CreateAlbum(ctx context.Context, arg db.CreateAlbumParams) (db.Album, error)
	UpdateAlbum(ctx context.Context, arg db.UpdateAlbumParams) (db.Album, error)
	DeleteAlbum(ctx context.Context, id int32) error
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

func (s *albumStore) GetAlbumsOptions(ctx context.Context, familyID int32, groupID int32) ([]db.GetAlbumsOptionsRow, error) {
	return wrapErr(s.queries.GetAlbumsOptions(ctx, db.GetAlbumsOptionsParams{
		FamilyID: familyID,
		GroupID:  groupID,
	}))
}

func (s *albumStore) FindAlbumsByFamilyID(ctx context.Context, familyID int32) ([]db.Album, error) {
	return wrapErr(s.queries.FindAlbumsByFamilyID(ctx, familyID))
}

func (s *albumStore) FindAlbumByIDAndFamilyID(ctx context.Context, id int32, familyID int32) (db.Album, error) {
	return wrapErr(s.queries.FindAlbumByIDAndFamilyID(ctx, db.FindAlbumByIDAndFamilyIDParams{
		ID:       id,
		FamilyID: familyID,
	}))
}

func (s *albumStore) CreateAlbum(ctx context.Context, arg db.CreateAlbumParams) (db.Album, error) {
	return wrapErr(s.queries.CreateAlbum(ctx, arg))
}

func (s *albumStore) UpdateAlbum(ctx context.Context, arg db.UpdateAlbumParams) (db.Album, error) {
	return wrapErr(s.queries.UpdateAlbum(ctx, arg))
}

func (s *albumStore) DeleteAlbum(ctx context.Context, id int32) error {
	return s.queries.DeleteAlbum(ctx, id)
}
