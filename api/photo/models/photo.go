package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Photo struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	ThumbnailUrl string             `json:"thumbnail_url" bson:"thumbnail_url"`
	OriginalUrl  string             `json:"original_url" bson:"original_url"`
	// LiveUrl         string             `json:"live_url,omitempty" bson:"live_url,omitempty"`
	// OriginalLiveUrl string             `json:"original_live_url" bson:"original_live_url"`
	FileName  string    `json:"file_name" bson:"file_name"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	CreatedBy string    `json:"created_by" bson:"created_by"`
	UpdatedBy string    `json:"updated_by" bson:"updated_by"`
}
