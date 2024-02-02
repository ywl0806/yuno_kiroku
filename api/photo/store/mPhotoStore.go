package store

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type MPhotoStore struct {
	collection *mongo.Collection
}

func NewMPhotoStore(db *mongo.Database) *MPhotoStore {
	return &MPhotoStore{
		collection: db.Collection("photos"),
	}
}

func (s *MPhotoStore) FindPictures() {}
