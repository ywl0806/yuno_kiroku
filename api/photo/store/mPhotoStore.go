package store

import (
	"context"
	"log"
	"time"

	"github.com/ywl0806/yuno_kiroku/api/photo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MPhotoStore represents a store for managing photos in MongoDB.
type MPhotoStore struct {
	collection *mongo.Collection
}

// NewMPhotoStore creates a new instance of MPhotoStore.
func NewMPhotoStore(db *mongo.Database) *MPhotoStore {
	return &MPhotoStore{
		collection: db.Collection("photos"),
	}
}

// FindPictures retrieves all photos from the store.
func (s *MPhotoStore) FindPictures(params *FindPictureParams) ([]models.Photo, error) {
	ctx := context.Background()
	limit := int64(*params.Limit)
	skip := int64(*params.Skip)
	options := options.FindOptions{Limit: &limit, Skip: &skip}
	cursor, err := s.collection.Find(ctx, bson.D{{}}, &options)
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

// CreatePicture creates a new photo in the store.
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
