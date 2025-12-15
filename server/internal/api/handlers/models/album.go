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
