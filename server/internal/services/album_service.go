package services

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/apperr"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store"
)

type GroupPermission struct {
	GroupID    int32
	Permission string
}

type AlbumService struct {
	albumStore           store.AlbumStore
	albumGroupPermStore  store.AlbumGroupPermissionStore
}

func NewAlbumService(albumStore store.AlbumStore, albumGroupPermStore store.AlbumGroupPermissionStore) *AlbumService {
	return &AlbumService{albumStore: albumStore, albumGroupPermStore: albumGroupPermStore}
}

func (s *AlbumService) GetAlbumsForWrite(ctx context.Context, familyId int32, groupId int32) ([]db.Album, error) {
	albums, err := s.albumStore.FindAlbumsForWrite(ctx, db.FindAlbumsForWriteParams{
		FamilyID: familyId,
		GroupID:  groupId,
	})
	if err != nil {
		return nil, err
	}
	return albums, nil
}

func (s *AlbumService) GetAlbumsOptions(ctx context.Context, familyID int32, groupID int32) ([]db.GetAlbumsOptionsRow, error) {
	return s.albumStore.GetAlbumsOptions(ctx, familyID, groupID)
}

func (s *AlbumService) GetAllAlbums(ctx context.Context, familyID int32) ([]db.Album, error) {
	return s.albumStore.FindAlbumsByFamilyID(ctx, familyID)
}

func (s *AlbumService) GetAlbumWithPermissions(ctx context.Context, albumID int32, familyID int32) (db.Album, []db.AlbumGroupsPermission, error) {
	album, err := s.albumStore.FindAlbumByID(ctx, albumID)
	if err != nil {
		return db.Album{}, nil, err
	}
	if album.FamilyID != familyID {
		return db.Album{}, nil, apperr.NewForbiddenError("error.forbidden", nil)
	}
	perms, err := s.albumGroupPermStore.GetAlbumGroupPermissions(ctx, albumID)
	if err != nil {
		return db.Album{}, nil, err
	}
	return album, perms, nil
}

func (s *AlbumService) CreateAlbum(ctx context.Context, familyID int32, name string, perms []GroupPermission) (db.Album, error) {
	album, err := s.albumStore.CreateAlbum(ctx, db.CreateAlbumParams{
		FamilyID: familyID,
		Name:     name,
	})
	if err != nil {
		return db.Album{}, err
	}
	for _, p := range perms {
		if err := s.albumGroupPermStore.InsertAlbumGroupPermission(ctx, db.InsertAlbumGroupPermissionParams{
			AlbumID:    album.ID,
			GroupID:    p.GroupID,
			Permission: p.Permission,
		}); err != nil {
			return db.Album{}, err
		}
	}
	return album, nil
}

func (s *AlbumService) UpdateAlbum(ctx context.Context, albumID int32, familyID int32, name string, perms []GroupPermission) (db.Album, error) {
	album, err := s.albumStore.FindAlbumByID(ctx, albumID)
	if err != nil {
		return db.Album{}, err
	}
	if album.FamilyID != familyID {
		return db.Album{}, apperr.NewForbiddenError("error.forbidden", nil)
	}
	updated, err := s.albumStore.UpdateAlbum(ctx, db.UpdateAlbumParams{
		Name: name,
		ID:   albumID,
	})
	if err != nil {
		return db.Album{}, err
	}
	if err := s.albumGroupPermStore.DeleteAlbumGroupPermissionsByAlbumID(ctx, albumID); err != nil {
		return db.Album{}, err
	}
	for _, p := range perms {
		if err := s.albumGroupPermStore.InsertAlbumGroupPermission(ctx, db.InsertAlbumGroupPermissionParams{
			AlbumID:    albumID,
			GroupID:    p.GroupID,
			Permission: p.Permission,
		}); err != nil {
			return db.Album{}, err
		}
	}
	return updated, nil
}

func (s *AlbumService) DeleteAlbum(ctx context.Context, albumID int32, familyID int32) error {
	album, err := s.albumStore.FindAlbumByID(ctx, albumID)
	if err != nil {
		return err
	}
	if album.FamilyID != familyID {
		return apperr.NewForbiddenError("error.forbidden", nil)
	}
	if err := s.albumGroupPermStore.DeleteAlbumGroupPermissionsByAlbumID(ctx, albumID); err != nil {
		return err
	}
	return s.albumStore.DeleteAlbum(ctx, albumID)
}
