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

type AlbumService interface {
	GetAlbumsForWrite(ctx context.Context, familyId int32, groupId int32) ([]db.Album, error)
	GetAlbumsOptions(ctx context.Context, familyID int32, groupID int32) ([]db.GetAlbumsOptionsRow, error)
	GetAllAlbums(ctx context.Context, familyID int32) ([]db.Album, error)
	GetAlbumWithPermissions(ctx context.Context, albumID int32, familyID int32) (db.Album, []db.AlbumGroupsPermission, error)
	CreateAlbum(ctx context.Context, familyID int32, name string, perms []GroupPermission) (db.Album, error)
	UpdateAlbum(ctx context.Context, albumID int32, familyID int32, name string, perms []GroupPermission) (db.Album, error)
	DeleteAlbum(ctx context.Context, albumID int32, familyID int32) error
}

type albumService struct {
	albumStore          store.AlbumStore
	albumGroupPermStore store.AlbumGroupPermissionStore
	transactor          store.Transactor
}

func NewAlbumService(albumStore store.AlbumStore, albumGroupPermStore store.AlbumGroupPermissionStore, transactor store.Transactor) AlbumService {
	return &albumService{albumStore: albumStore, albumGroupPermStore: albumGroupPermStore, transactor: transactor}
}

func (s *albumService) GetAlbumsForWrite(ctx context.Context, familyId int32, groupId int32) ([]db.Album, error) {
	albums, err := s.albumStore.FindAlbumsForWrite(ctx, db.FindAlbumsForWriteParams{
		FamilyID: familyId,
		GroupID:  groupId,
	})
	if err != nil {
		return nil, err
	}
	return albums, nil
}

func (s *albumService) GetAlbumsOptions(ctx context.Context, familyID int32, groupID int32) ([]db.GetAlbumsOptionsRow, error) {
	return s.albumStore.GetAlbumsOptions(ctx, familyID, groupID)
}

func (s *albumService) GetAllAlbums(ctx context.Context, familyID int32) ([]db.Album, error) {
	return s.albumStore.FindAlbumsByFamilyID(ctx, familyID)
}

func (s *albumService) GetAlbumWithPermissions(ctx context.Context, albumID int32, familyID int32) (db.Album, []db.AlbumGroupsPermission, error) {
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

func (s *albumService) CreateAlbum(ctx context.Context, familyID int32, name string, perms []GroupPermission) (db.Album, error) {
	var result db.Album
	err := s.transactor.Transact(ctx, func(tx store.TxStore) error {
		album, err := tx.Album().CreateAlbum(ctx, db.CreateAlbumParams{
			FamilyID: familyID,
			Name:     name,
		})
		if err != nil {
			return err
		}
		for _, p := range perms {
			if err := tx.AlbumGroupPermission().InsertAlbumGroupPermission(ctx, db.InsertAlbumGroupPermissionParams{
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

func (s *albumService) UpdateAlbum(ctx context.Context, albumID int32, familyID int32, name string, perms []GroupPermission) (db.Album, error) {
	// 트랜잭션 외부에서 권한 검증 (조회만)
	album, err := s.albumStore.FindAlbumByID(ctx, albumID)
	if err != nil {
		return db.Album{}, err
	}
	if album.FamilyID != familyID {
		return db.Album{}, apperr.NewForbiddenError("error.forbidden", nil)
	}

	var result db.Album
	err = s.transactor.Transact(ctx, func(tx store.TxStore) error {
		updated, err := tx.Album().UpdateAlbum(ctx, db.UpdateAlbumParams{
			Name: name,
			ID:   albumID,
		})
		if err != nil {
			return err
		}
		if err := tx.AlbumGroupPermission().DeleteAlbumGroupPermissionsByAlbumID(ctx, albumID); err != nil {
			return err
		}
		for _, p := range perms {
			if err := tx.AlbumGroupPermission().InsertAlbumGroupPermission(ctx, db.InsertAlbumGroupPermissionParams{
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

func (s *albumService) DeleteAlbum(ctx context.Context, albumID int32, familyID int32) error {
	album, err := s.albumStore.FindAlbumByID(ctx, albumID)
	if err != nil {
		return err
	}
	if album.FamilyID != familyID {
		return apperr.NewForbiddenError("error.forbidden", nil)
	}
	return s.transactor.Transact(ctx, func(tx store.TxStore) error {
		if err := tx.AlbumGroupPermission().DeleteAlbumGroupPermissionsByAlbumID(ctx, albumID); err != nil {
			return err
		}
		return tx.Album().DeleteAlbum(ctx, albumID)
	})
}
