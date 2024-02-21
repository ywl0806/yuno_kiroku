package store

import (
	"context"
	"log"
	"time"

	"github.com/ywl0806/yuno_kiroku/api/photo/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (s *MPhotoStore) FindPictures() {
	ctx := context.Background()

	s.collection.Find(ctx, nil)
}

func (s *MPhotoStore) CreatePicture(photo models.Photo) (models.Photo, error) {
	ctx := context.Background()

	createdAt := time.Now()
	photo.CreatedAt = createdAt
	photo.UpdatedAt = createdAt

	result, err := s.collection.InsertOne(ctx, photo)
	if err != nil {
		log.Println("insert photo error: ", err)
		return models.Photo{}, err
	}

	photo.ID = result.InsertedID.(primitive.ObjectID)

	return photo, nil
}
