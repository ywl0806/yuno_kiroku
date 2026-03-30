package storage

import "time"

// Storage서비스의 인터페이스
type StorageService interface {
	SaveFile(file []byte, filePath string, fileName string) (string, error)
	GetFile(key string) ([]byte, error)
	GeneratePresignedPutURL(key string, contentType string, expiresIn time.Duration) (string, error)
}
