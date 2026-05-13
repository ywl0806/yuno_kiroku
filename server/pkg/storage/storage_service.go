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
	// DownloadToFile S3 객체를 로컬 파일로 스트리밍 다운로드 (대용량 파일용)
	DownloadToFile(ctx context.Context, key string, destPath string) error
	// UploadFromFile 로컬 파일을 S3에 스트리밍 업로드 (대용량 파일용)
	UploadFromFile(ctx context.Context, key string, contentType string, srcPath string) error
}
