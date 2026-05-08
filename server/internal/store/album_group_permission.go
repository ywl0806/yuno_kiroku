package store

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

// AlbumGroupPermissionStore 앨범-그룹 권한 데이터 접근 인터페이스
type AlbumGroupPermissionStore interface {
	GetAlbumGroupPermissions(ctx context.Context, albumID int32, familyID int32) ([]db.AlbumGroupsPermission, error)
	InsertAlbumGroupPermission(ctx context.Context, arg db.InsertAlbumGroupPermissionParams) error
	DeleteAlbumGroupPermissionsByAlbumID(ctx context.Context, albumID int32) error
}

type albumGroupPermissionStore struct {
	queries *db.Queries
}

func NewAlbumGroupPermissionStore(queries *db.Queries) AlbumGroupPermissionStore {
	return &albumGroupPermissionStore{queries: queries}
}

func (s *albumGroupPermissionStore) GetAlbumGroupPermissions(ctx context.Context, albumID int32, familyID int32) ([]db.AlbumGroupsPermission, error) {
	return wrapErr(s.queries.GetAlbumGroupPermissions(ctx, db.GetAlbumGroupPermissionsParams{
		AlbumID:  albumID,
		FamilyID: familyID,
	}))
}

func (s *albumGroupPermissionStore) InsertAlbumGroupPermission(ctx context.Context, arg db.InsertAlbumGroupPermissionParams) error {
	return s.queries.InsertAlbumGroupPermission(ctx, arg)
}

func (s *albumGroupPermissionStore) DeleteAlbumGroupPermissionsByAlbumID(ctx context.Context, albumID int32) error {
	return s.queries.DeleteAlbumGroupPermissionsByAlbumID(ctx, albumID)
}
