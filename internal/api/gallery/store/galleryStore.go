package store

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GalleryStore struct {
	ID          primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	Members     []primitive.ObjectID `json:"members" bson:"members"`
	Name        string               `json:"name" bson:"name"`
	Description string               `json:"description" bson:"description"`
	CreatedAt   time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at" bson:"updated_at"`
	CreatedBy   string               `json:"created_by" bson:"created_by"`
	UpdatedBy   string               `json:"updated_by" bson:"updated_by"`
}
