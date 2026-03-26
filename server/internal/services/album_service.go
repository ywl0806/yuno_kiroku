package services

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store"
)

type AlbumService struct {
	albumStore store.AlbumStore
}

func NewAlbumService(albumStore store.AlbumStore) *AlbumService {
	return &AlbumService{albumStore: albumStore}
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
