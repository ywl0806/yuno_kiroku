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
	albumStore          store.AlbumStore
	albumGroupPermStore store.AlbumGroupPermissionStore
	transactor          store.Transactor
}

func NewAlbumService(albumStore store.AlbumStore, albumGroupPermStore store.AlbumGroupPermissionStore, transactor store.Transactor) *AlbumService {
	return &AlbumService{albumStore: albumStore, albumGroupPermStore: albumGroupPermStore, transactor: transactor}
}

func (s *AlbumService) GetAlbumsForWrite(ctx context.Context, familyId string, groupId int32) ([]db.Album, error) {
	albums, err := s.albumStore.FindAlbumsForWrite(ctx, db.FindAlbumsForWriteParams{
		FamilyID: familyId,
		GroupID:  groupId,
	})
	if err != nil {
		return nil, err
	}
	return albums, nil
}

func (s *AlbumService) GetAlbumsOptions(ctx context.Context, familyID string, groupID int32) ([]db.GetAlbumsOptionsRow, error) {
	return s.albumStore.GetAlbumsOptions(ctx, familyID, groupID)
}

func (s *AlbumService) GetAllAlbums(ctx context.Context, familyID string) ([]db.Album, error) {
	return s.albumStore.FindAlbumsByFamilyID(ctx, familyID)
}

func (s *AlbumService) GetAlbum(ctx context.Context, albumID string, familyID string) (db.Album, error) {
	album, err := s.albumStore.FindAlbumByIDAndFamilyID(ctx, albumID, familyID)
	if err != nil {
		return db.Album{}, err
	}
	return album, nil
}

func (s *AlbumService) GetAlbumPermissions(ctx context.Context, albumID string, familyID string) ([]db.AlbumGroupsPermission, error) {
	perms, err := s.albumGroupPermStore.GetAlbumGroupPermissions(ctx, albumID, familyID)
	if err != nil {
		return nil, err
	}
	return perms, nil
}

func (s *AlbumService) GetAlbumWithPermissions(ctx context.Context, albumID string, familyID string) (db.Album, []db.AlbumGroupsPermission, error) {
	album, err := s.albumStore.FindAlbumByIDAndFamilyID(ctx, albumID, familyID)
	if err != nil {
		return db.Album{}, nil, err
	}
	perms, err := s.albumGroupPermStore.GetAlbumGroupPermissions(ctx, albumID, familyID)
	if err != nil {
		return db.Album{}, nil, err
	}
	return album, perms, nil
}

func (s *AlbumService) CreateAlbum(ctx context.Context, familyID string, name string, isCommon bool, perms []GroupPermission) (db.Album, error) {
	var result db.Album
	err := s.transactor.Transact(ctx, func(tx *store.Store) error {
		album, err := tx.Album.CreateAlbum(ctx, db.CreateAlbumParams{
			FamilyID: familyID,
			Name:     name,
			IsCommon: isCommon,
		})
		if err != nil {
			return err
		}
		for _, p := range perms {
			if err := tx.AlbumGroupPermission.InsertAlbumGroupPermission(ctx, db.InsertAlbumGroupPermissionParams{
				AlbumID:    album.ID,
				GroupID:    p.GroupID,
				Permission: p.Permission,
			}); err != nil {
				return err
			}
		}
		result = album
		return nil
	})
	return result, err
}

func (s *AlbumService) UpdateAlbum(ctx context.Context, albumID string, familyID string, name string, perms []GroupPermission) (db.Album, error) {
	album, err := s.albumStore.FindAlbumByIDAndFamilyID(ctx, albumID, familyID)
	if err != nil {
		return db.Album{}, err
	}
	if album.FamilyID != familyID {
		return db.Album{}, apperr.NewForbiddenError("error.forbidden", nil)
	}
	if album.IsCommon {
		return db.Album{}, apperr.NewForbiddenError("error.album.common_not_editable", nil)
	}

	var result db.Album
	err = s.transactor.Transact(ctx, func(tx *store.Store) error {
		updated, err := tx.Album.UpdateAlbum(ctx, db.UpdateAlbumParams{
			Name: name,
			ID:   albumID,
		})
		if err != nil {
			return err
		}
		if err := tx.AlbumGroupPermission.DeleteAlbumGroupPermissionsByAlbumID(ctx, albumID); err != nil {
			return err
		}
		for _, p := range perms {
			if err := tx.AlbumGroupPermission.InsertAlbumGroupPermission(ctx, db.InsertAlbumGroupPermissionParams{
				AlbumID:    albumID,
				GroupID:    p.GroupID,
				Permission: p.Permission,
			}); err != nil {
				return err
			}
		}
		result = updated
		return nil
	})
	return result, err
}

func (s *AlbumService) DeleteAlbum(ctx context.Context, albumID string, familyID string) error {
	album, err := s.albumStore.FindAlbumByIDAndFamilyID(ctx, albumID, familyID)
	if err != nil {
		return err
	}

	if album.IsCommon {
		return apperr.NewForbiddenError("error.album.common_not_deletable", nil)
	}
	return s.transactor.Transact(ctx, func(tx *store.Store) error {
		if err := tx.AlbumGroupPermission.DeleteAlbumGroupPermissionsByAlbumID(ctx, albumID); err != nil {
			return err
		}
		return tx.Album.DeleteAlbum(ctx, albumID)
	})
}
