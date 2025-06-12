package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SocialAccount struct {
	Provider   string `bson:"provider" json:"provider"`
	ProviderID string `bson:"provider_id" json:"provider_id"`
}
type User struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	GroupID        primitive.ObjectID `json:"group_id" bson:"group_id,omitempty"`
	Email          *string            `json:"email" bson:"email" validate:"required,email"`
	Name           string             `json:"name" bson:"name" validate:"required"`
	Password       *string            `json:"password" bson:"password" validate:"required,min=6"`
	IsActive       bool               `json:"is_active" bson:"is_active"`
	LastLogin      *time.Time         `json:"last_login" bson:"last_login"`
	CreatedAt      time.Time          `json:"create_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
	SocialAccounts []SocialAccount    `bson:"social_accounts,omitempty" json:"social_accounts,omitempty"`
	RefreshToken   *string            `json:"refresh_token" bson:"refresh_token,omitempty"`
}
