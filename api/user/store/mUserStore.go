package store

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

type MUserStore struct {
	collection *mongo.Collection
}

func NewMUserStore(db *mongo.Database) *MUserStore {
	return &MUserStore{
		collection: db.Collection("users"),
	}
}

func (s *MUserStore) FindUsers() {
	ctx := context.Background()
	users, _ := s.collection.Find(ctx, nil)
	fmt.Println(users)
}
func (s *MUserStore) FindUserById() {}
func (s *MUserStore) CreateUser() {
	ctx := context.Background()
	result, err := s.collection.InsertOne(ctx, map[string]interface{}{
		"name": "hoge",
	})

	fmt.Println(result)

	if err != nil {
		log.Fatal(err)
	}
}
