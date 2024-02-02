package storage

import "mime/multipart"

// Storage서비스의 인터페이스
type StorageService interface {
	SaveFile(file *multipart.FileHeader, filePath string) (string, error)
}
