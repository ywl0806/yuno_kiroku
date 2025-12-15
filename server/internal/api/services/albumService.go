package services

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

type AlbumService struct {
	queries *db.Queries
}

func NewAlbumService(queries *db.Queries) *AlbumService {
	return &AlbumService{queries: queries}
}

func (s *AlbumService) GetAlbumsForWrite(ctx context.Context, groupId int32, clanGroupId int32) ([]db.Album, error) {
	albums, err := s.queries.FindAlbumsForWrite(ctx, db.FindAlbumsForWriteParams{
		GroupID:     groupId,
		ClanGroupID: clanGroupId,
	})
	if err != nil {
		return nil, err
	}
	return albums, nil
}
