package store

import (
	"context"
	"log"
	"time"

	"github.com/ywl0806/yuno_kiroku/api/photo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FindPicturesGroupByDate retrieves all photos from the store grouped by date.
func (s *PhotoStore) FindPicturesGroupByDate(from, to time.Time) ([]models.PhotoGroup, error) {
	ctx := context.Background()

	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "year", Value: bson.D{{Key: "$year", Value: "$photo_created_at"}}},
			{Key: "month", Value: bson.D{{Key: "$month", Value: "$photo_created_at"}}},
			{Key: "thumbnail_url", Value: 1},
			{Key: "original_url", Value: 1},
			{Key: "live_url", Value: 1},
			{Key: "original_live_url", Value: 1},
			{Key: "file_name", Value: 1},
			{Key: "photo_created_at", Value: 1},
			{Key: "created_at", Value: 1},
			{Key: "updated_at", Value: 1},
			{Key: "created_by", Value: 1},
			{Key: "updated_by", Value: 1},
		}},
	}

	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "photo_created_at", Value: bson.D{
				{Key: "$gte", Value: from},
				{Key: "$lt", Value: to},
			}},
		}},
	}

	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "year", Value: "$year"},
				{Key: "month", Value: "$month"},
			}},
			{Key: "photos", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		}},
	}

	// Create a pipeline to group photos by date.
	pipeline := mongo.Pipeline{projectStage, matchStage, groupStage}

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("aggregate photo error: ", err)
		return nil, err
	}

	var photoGroups []models.PhotoGroup
	if err = cursor.All(ctx, &photoGroups); err != nil {
		log.Println("cursor all error: ", err)
		return nil, err
	}

	return photoGroups, nil
}
