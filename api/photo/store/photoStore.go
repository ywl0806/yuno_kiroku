package store

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FindPictureParams struct {
	// Limit is the maximum number of photos to return.
	Limit *int
	// Skip is the number of photos to skip.
	Skip *int
}

// PhotoStore represents a store for managing photos in MongoDB.
type PhotoStore struct {
	collection *mongo.Collection
}

// NewPhotoStore creates a new instance of PhotoStore.
func NewPhotoStore(db *mongo.Database) *PhotoStore {
	col := db.Collection("photos")

	col.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{{Key: "photo_created_at", Value: 1}},
		},
	)

	return &PhotoStore{
		collection: db.Collection("photos"),
	}
}
