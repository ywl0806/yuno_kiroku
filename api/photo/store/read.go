package store

import (
	"context"
	"log"

	"github.com/ywl0806/yuno_kiroku/api/photo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindPictures retrieves all photos from the store.
func (s *PhotoStore) FindPictures(params *FindPictureParams) ([]models.Photo, error) {
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

func (s *PhotoStore) FindPhotosRange() ([]models.PhotoRange, error) {
	ctx := context.Background()

	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "year", Value: bson.D{{Key: "$year", Value: "$photo_created_at"}}},
			{Key: "month", Value: bson.D{{Key: "$month", Value: "$photo_created_at"}}},
		}},
	}

	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{
				Key: "photo_range",
				Value: bson.D{{
					Key: "$addToSet",
					Value: bson.D{
						{Key: "year", Value: "$year"},
						{Key: "month", Value: "$month"},
					},
				}},
			},
		},
		}}

	// Create a pipeline to group photos by date.
	pipeline := mongo.Pipeline{projectStage, groupStage}

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("aggregate photo error: ", err)
		return nil, err
	}

	defer cursor.Close(ctx)
	cursor.Next(ctx)

	var photoRange []models.PhotoRanges
	if err = cursor.All(ctx, &photoRange); err != nil {
		log.Println("cursor all error: ", err)
		return nil, err
	}

	return photoRange[0].PhotoRange, nil
}
func (s *PhotoStore) FindOnePhoto(opt *options.FindOneOptions) (models.Photo, error) {
	ctx := context.Background()

	cursor := s.collection.FindOne(ctx, bson.D{{}}, opt)
	if cursor.Err() != nil {
		log.Println("find photo error: ", cursor.Err())
		return models.Photo{}, cursor.Err()
	}

	var photo models.Photo
	if err := cursor.Decode(&photo); err != nil {
		log.Println("cursor all error: ", err)
		return models.Photo{}, err
	}

	return photo, nil
}
