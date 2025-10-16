package db

import (
	"context"
	"log"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(viper.GetString("MONGO_URI")))

	if err != nil {
		log.Fatal("DB connection failed : ", err)
	}

	// err = client.Ping(ctx, readpref.Primary())

	// if err != nil {
	// 	log.Fatal("DB ping failed : ", err)
	// }

	return client
}
