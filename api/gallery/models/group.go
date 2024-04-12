package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Group struct {
	ID          primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	Members     []primitive.ObjectID `json:"members" bson:"members"`
	Admins      []primitive.ObjectID `json:"admins" bson:"admins"`
	Albums      []primitive.ObjectID `json:"albums" bson:"albums"`
	Name        string               `json:"name" bson:"name"`
	Description string               `json:"description" bson:"description"`

	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt primitive.DateTime `json:"updated_at" bson:"updated_at"`
	CreatedBy primitive.ObjectID `json:"created_by" bson:"created_by"`
	UpdatedBy primitive.ObjectID `json:"updated_by" bson:"updated_by"`
}
