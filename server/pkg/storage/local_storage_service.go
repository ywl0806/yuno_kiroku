package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)


// 로컬 스토리지 서비스의 구현체
type LocalStorageService struct {
	rootDir string
}

func NewLocalStorageService(rootDir string) *LocalStorageService {
	return &LocalStorageService{
		rootDir: rootDir,
	}
}

func (s *LocalStorageService) GetFile(key string) ([]byte, error) {
	return os.ReadFile(key)
}

func (s *LocalStorageService) GeneratePresignedPutURL(key string, contentType string, expiresIn time.Duration) (string, error) {
	return "", fmt.Errorf("local storage does not support presigned URLs")
}

func (s *LocalStorageService) SaveFile(file []byte, filePath string, fileName string) (string, error) {

	dirPath := filepath.Join("uploads", s.rootDir, filePath)

	err := os.MkdirAll(dirPath, os.ModePerm)
	path := filepath.Join(dirPath, fileName)

	if err != nil {
		return "", err
	}

	dst, err := os.Create(path)

	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := dst.Write(file); err != nil {
		return "", err
	}

	return path, err
}
