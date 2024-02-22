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

func (s *MPhotoStore) FindPictures() ([]models.Photo, error) {
	ctx := context.Background()

	cursor, err := s.collection.Find(ctx, nil)
	if err != nil {
		log.Println("find photo error: ", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var photos []models.Photo
	if err = cursor.All(ctx, &photos); err != nil {
		log.Println("cursor all error: ", err)
		return nil, err
	}

	return photos, nil
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
