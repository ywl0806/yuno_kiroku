package store

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

// AlbumGroupPermissionStore 앨범-그룹 권한 데이터 접근 인터페이스
type AlbumGroupPermissionStore interface {
	GetAlbumGroupPermissions(ctx context.Context, albumID string, familyID string) ([]db.AlbumGroupsPermission, error)
	InsertAlbumGroupPermission(ctx context.Context, arg db.InsertAlbumGroupPermissionParams) error
	DeleteAlbumGroupPermissionsByAlbumID(ctx context.Context, albumID string) error
	CheckUserHasPermissionForAlbum(ctx context.Context, albumID string, familyID string, userID string) (bool, error)
	GetWritableAlbumIDsByGroupID(ctx context.Context, groupID int32) ([]string, error)
}

type albumGroupPermissionStore struct {
	queries *db.Queries
}

func NewAlbumGroupPermissionStore(queries *db.Queries) AlbumGroupPermissionStore {
	return &albumGroupPermissionStore{queries: queries}
}

func (s *albumGroupPermissionStore) GetAlbumGroupPermissions(ctx context.Context, albumID string, familyID string) ([]db.AlbumGroupsPermission, error) {
	return wrapErr(s.queries.GetAlbumGroupPermissions(ctx, db.GetAlbumGroupPermissionsParams{
		AlbumID:  albumID,
		FamilyID: familyID,
	}))
}

func (s *albumGroupPermissionStore) InsertAlbumGroupPermission(ctx context.Context, arg db.InsertAlbumGroupPermissionParams) error {
	return mapDBError(s.queries.InsertAlbumGroupPermission(ctx, arg))
}

func (s *albumGroupPermissionStore) DeleteAlbumGroupPermissionsByAlbumID(ctx context.Context, albumID string) error {
	return mapDBError(s.queries.DeleteAlbumGroupPermissionsByAlbumID(ctx, albumID))
}

func (s *albumGroupPermissionStore) CheckUserHasPermissionForAlbum(ctx context.Context, albumID string, userID string, permission string) (bool, error) {
	return wrapErr(s.queries.CheckUserHasPermissionForAlbum(ctx, db.CheckUserHasPermissionForAlbumParams{
		AlbumID:    albumID,
		UserID:     userID,
		Permission: permission,
	}))
}

func (s *albumGroupPermissionStore) GetWritableAlbumIDsByGroupID(ctx context.Context, groupID int32) ([]string, error) {
	return wrapErr(s.queries.GetWritableAlbumIDsByGroupID(ctx, groupID))
}
