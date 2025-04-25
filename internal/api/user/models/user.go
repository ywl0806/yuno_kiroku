package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"user_name"`
	Name      string             `json:"name" bson:"name"`
	Password  string             `json:"password" bson:"password"`
	CreatedAt time.Time          `json:"create_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
