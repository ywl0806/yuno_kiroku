package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Photo struct {
	ID              primitive.ObjectID `bson:"_id"`
	ThumbnailUrl    string             `bson:"thumbnail_url"`
	LiveUrl         string             `bson:"live_url,omitempty"`
	OriginalUrl     string             `bson:"original_url"`
	OriginalLiveUrl string             `bson:"original_live_url"`
	FileName        string             `bson:"file_name"`
	CreatedAt       time.Time          `bson:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at"`
}
