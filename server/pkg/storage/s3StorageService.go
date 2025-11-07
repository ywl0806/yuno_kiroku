package storage

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// s3 storage service

type S3StorageService struct {
	bucketName string
	client     *s3.Client
}

func NewS3StorageService(bucketName string) *S3StorageService {

	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
		panic(err)
	}

	client := s3.NewFromConfig(cfg)
	return &S3StorageService{bucketName: bucketName, client: client}
}

func (s *S3StorageService) SaveFile(file io.Reader, filePath string, fileName string) (string, error) {
	ctx := context.Background()

	fileKey := fmt.Sprintf("%s/%s", filePath, fileName)

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileKey),
		Body:   file,
	})

	if err != nil {
		return "", err
	}
	return fileKey, nil
}
