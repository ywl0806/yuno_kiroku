package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Username  string             `bson:"user_name"`
	Name      string             `bson:"name"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type UserCred struct {
	ID       primitive.ObjectID `bson:"_id"`
	Password string             `bson:"password"`

	CreatedAt time.Time `bson:"create_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
