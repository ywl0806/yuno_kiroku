package store

import (
	"context"
	"time"

	"github.com/ywl0806/yuno_kiroku/internal/api/services/group/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GroupStore struct {
	collection *mongo.Collection
}

func NewGroupStore(db *mongo.Database) *GroupStore {
	return &GroupStore{
		collection: db.Collection("groups"),
	}
}

func (s *GroupStore) CreateGroup(ctx context.Context, name string) (models.Group, error) {

	group := models.Group{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := s.collection.InsertOne(ctx, group)
	if err != nil {
		return models.Group{}, err
	}

	group.ID = result.InsertedID.(primitive.ObjectID)

	return group, nil
}
