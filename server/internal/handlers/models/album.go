package models

import (
	"time"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

type AlbumResponse struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewAlbumResponse(album *db.Album) *AlbumResponse {
	return &AlbumResponse{
		ID:        album.ID,
		Name:      album.Name,
		CreatedAt: album.CreatedAt,
		UpdatedAt: album.UpdatedAt,
	}
}

type AlbumOptionResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func NewAlbumOptionResponse(albumOption *db.GetAlbumsOptionsRow) *AlbumOptionResponse {
	return &AlbumOptionResponse{
		ID:   albumOption.ID,
		Name: albumOption.Name,
	}
}

func NewAlbumOptionResponses(albumOptions []db.GetAlbumsOptionsRow) *[]AlbumOptionResponse {
	albumOptionResponses := make([]AlbumOptionResponse, len(albumOptions))
	for i, albumOption := range albumOptions {
		albumOptionResponses[i] = *NewAlbumOptionResponse(&albumOption)
	}
	return &albumOptionResponses
}
