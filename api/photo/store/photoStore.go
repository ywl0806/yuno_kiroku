package store

import "github.com/ywl0806/yuno_kiroku/api/photo/models"

type FindPictureParams struct {
	// Limit is the maximum number of photos to return.
	Limit *int
	// Skip is the number of photos to skip.
	Skip *int
}
type PhotoStore interface {
	FindPictures(params *FindPictureParams) ([]models.Photo, error)
	CreatePicture(photo models.Photo) (models.Photo, error)
}
