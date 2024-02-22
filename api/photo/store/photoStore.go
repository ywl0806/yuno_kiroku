package store

import "github.com/ywl0806/yuno_kiroku/api/photo/models"

type PhotoStore interface {
	FindPictures() ([]models.Photo, error)
	CreatePicture(photo models.Photo) (models.Photo, error)
}
