package store

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type GroupStore struct {
	collection *mongo.Collection
}

func NewGroupStore(db *mongo.Database) *GroupStore {
	col := db.Collection("clans")

	// Create indexes if needed
	// col.Indexes().CreateOne(
	// 	context.Background(),
	// 	mongo.IndexModel{
	// 		Keys: bson.D{{Key: "name", Value: 1}},
	// 	},
	// )

	return &GroupStore{
		collection: col,
	}
}
