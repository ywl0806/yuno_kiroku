package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
)

// s3 storage service

type S3StorageService struct {
	bucketName    string
	client        *s3.Client
	presignClient *s3.PresignClient
}

func NewS3StorageService(bucketName string) *S3StorageService {

	storageType := viper.GetString("STORAGE_TYPE")

	if storageType == "" {
		storageType = "s3"
	}

	var client *s3.Client
	var presignClient *s3.PresignClient

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

		// 외부 공개 호스트가 있으면 해당 엔드포인트로 presign 클라이언트 별도 생성
		// (서명의 host 헤더가 클라이언트가 실제 요청하는 호스트와 일치해야 함)
		if storageServiceHost := viper.GetString("STORAGE_SERVICE_HOST"); storageServiceHost != "" {
			externalClient := s3.NewFromConfig(cfg, func(o *s3.Options) {
				o.BaseEndpoint = aws.String(storageServiceHost)
				o.UsePathStyle = true
			})
			presignClient = s3.NewPresignClient(externalClient)
		}

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

	return &S3StorageService{bucketName: bucketName, client: client, presignClient: presignClient}
}

func (s *S3StorageService) GetFile(ctx context.Context, key string) ([]byte, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("파일 다운로드 실패: %w", err)
	}
	defer result.Body.Close()
	return io.ReadAll(result.Body)
}

func (s *S3StorageService) GeneratePresignedPutURL(ctx context.Context, key string, contentType string, expiresIn time.Duration) (string, error) {
	pc := s.presignClient
	if pc == nil {
		pc = s3.NewPresignClient(s.client)
	}

	req, err := pc.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(expiresIn))

	if err != nil {
		return "", fmt.Errorf("presigned URL 생성 실패: %w", err)
	}

	return req.URL, nil
}

func (s *S3StorageService) DeleteFile(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("파일 삭제 실패: %w", err)
	}
	return nil
}

func (s *S3StorageService) SaveFile(ctx context.Context, key string, file []byte) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(file),
	})
	if err != nil {
		return "", err
	}
	return key, nil
}

func (s *S3StorageService) DownloadToFile(ctx context.Context, key string, destPath string) error {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("S3 GetObject 실패 (key=%s): %w", key, err)
	}
	defer result.Body.Close()

	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("로컬 파일 생성 실패: %w", err)
	}
	defer f.Close()

	if _, err = io.Copy(f, result.Body); err != nil {
		return fmt.Errorf("스트리밍 다운로드 실패: %w", err)
	}
	return nil
}

func (s *S3StorageService) GetFileSize(ctx context.Context, key string) (int64, error) {
	result, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return 0, fmt.Errorf("HeadObject 실패: %w", err)
	}
	if result.ContentLength == nil {
		return 0, nil
	}
	return *result.ContentLength, nil
}

func (s *S3StorageService) UploadFromFile(ctx context.Context, key string, contentType string, srcPath string) error {
	f, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("로컬 파일 열기 실패: %w", err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("파일 정보 읽기 실패: %w", err)
	}

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucketName),
		Key:           aws.String(key),
		Body:          f,
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(fi.Size()),
	})
	return err
}
