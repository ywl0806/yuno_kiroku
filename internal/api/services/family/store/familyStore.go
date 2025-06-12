package store

import (
	"github.com/ywl0806/yuno_kiroku/internal/api/services/family/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type FamilyStore struct {
	collection *mongo.Collection
}

func NewFamilyStore(db *mongo.Database) *FamilyStore {
	col := db.Collection("families")

	// Create indexes if needed
	// col.Indexes().CreateOne(
	// 	context.Background(),
	// 	mongo.IndexModel{
	// 		Keys: bson.D{{Key: "name", Value: 1}},
	// 	},
	// )

	return &FamilyStore{
		collection: col,
	}
}

func (s *FamilyStore) CreateFamily() models.Family {
	return models.Family{
		// Initialize with default values if needed
	}
}
