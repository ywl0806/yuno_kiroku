package db

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	region string = "ap-northeast-1"
)

func NewClient(profile string) *dynamodb.Client {

	var ctx = context.Background()
	c, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))

	if err != nil {
		log.Fatal("DynamoDB connection error")
		return nil
	}
	client := dynamodb.NewFromConfig(c)
	return client
}
