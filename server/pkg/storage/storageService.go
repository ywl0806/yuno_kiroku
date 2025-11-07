package storage

import (
	"io"
)

// Storage서비스의 인터페이스
type StorageService interface {
	SaveFile(file io.Reader, filePath string, fileName string) (string, error)
}
