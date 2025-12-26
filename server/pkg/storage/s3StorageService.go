package storage

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
)

// s3 storage service

type S3StorageService struct {
	bucketName string
	client     *s3.Client
}

func NewS3StorageService(bucketName string) *S3StorageService {

	storageType := viper.GetString("STORAGE_TYPE")
	var client *s3.Client

	if storageType == "minio" {
		// MinIO를 위한 커스텀 설정
		creds := credentials.NewStaticCredentialsProvider(viper.GetString("MINIO_ROOT_USER"), viper.GetString("MINIO_ROOT_PASSWORD"), "")

		// MinIO를 위한 기본 config 로드
		cfg, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion("us-east-1"),
			config.WithCredentialsProvider(creds),
		)
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
			panic(err)
		}

		client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(viper.GetString("S3_ENDPOINT"))
			o.UsePathStyle = true
		})

	} else {
		// AWS S3를 위한 기본 설정
		cfg, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
			panic(err)
		}
		if region := viper.GetString("AWS_REGION"); region != "" {
			cfg.Region = region
		}
		client = s3.NewFromConfig(cfg)
	}

	return &S3StorageService{bucketName: bucketName, client: client}
}

func (s *S3StorageService) SaveFile(file []byte, filePath string, fileName string) (string, error) {
	ctx := context.Background()

	fileKey := fmt.Sprintf("%s/%s", filePath, fileName)

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileKey),
		Body:   bytes.NewReader(file),
	})

	if err != nil {
		return "", err
	}
	return fileKey, nil
}
