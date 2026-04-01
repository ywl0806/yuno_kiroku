package storage

import (
	"context"
	"time"
)

// Storage서비스의 인터페이스
type StorageService interface {
	SaveFile(ctx context.Context, file []byte, filePath string, fileName string) (string, error)
	GetFile(ctx context.Context, key string) ([]byte, error)
	GeneratePresignedPutURL(ctx context.Context, key string, contentType string, expiresIn time.Duration) (string, error)
}
