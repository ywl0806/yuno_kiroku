package store

import (
	"context"

	"cloud.google.com/go/firestore"
)

type UserStore struct {
	collection *firestore.CollectionRef
}

func NewUserStore(client *firestore.Client) *UserStore {
	return &UserStore{
		collection: client.Collection("users"),
	}
}

func (u *UserStore) Add() {
	ctx := context.Background()
	u.collection.Add(ctx, map[string]interface{}{
		"name": "hoge",
	})

}
