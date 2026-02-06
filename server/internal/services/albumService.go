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

func (s *AlbumService) GetAlbumsForWrite(ctx context.Context, groupId int32, clanGroupId int32) ([]db.Album, error) {
	albums, err := s.albumStore.FindAlbumsForWrite(ctx, db.FindAlbumsForWriteParams{
		GroupID:     groupId,
		ClanGroupID: clanGroupId,
	})
	if err != nil {
		return nil, err
	}
	return albums, nil
}
