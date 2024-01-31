package db

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

func New() *firestore.Client {
	os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8080")
	ctx := context.Background()
	conf := &firebase.Config{
		ProjectID:   "devlocal",
		DatabaseURL: "localhost:8080",
	}

	app, err := firebase.NewApp(ctx, conf)

	if err != nil {
		log.Fatal(err)
	}
	client, err := app.Firestore(ctx)

	if err != nil {
		log.Fatal(err)
	}
	return client
}
