package storage

import (
	"context"
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

func (s *LocalStorageService) GetFile(ctx context.Context, key string) ([]byte, error) {
	return os.ReadFile(key)
}

func (s *LocalStorageService) GeneratePresignedPutURL(ctx context.Context, key string, contentType string, expiresIn time.Duration) (string, error) {
	return "", fmt.Errorf("local storage does not support presigned URLs")
}

func (s *LocalStorageService) DeleteFile(ctx context.Context, key string) error {
	path := filepath.Join("uploads", s.rootDir, key)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("파일 삭제 실패: %w", err)
	}
	return nil
}

func (s *LocalStorageService) SaveFile(ctx context.Context, key string, file []byte) (string, error) {
	path := filepath.Join("uploads", s.rootDir, key)
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, file, 0644); err != nil {
		return "", err
	}
	return path, nil
}

func (s *LocalStorageService) DownloadToFile(ctx context.Context, key string, destPath string) error {
	data, err := os.ReadFile(key)
	if err != nil {
		return err
	}
	return os.WriteFile(destPath, data, 0644)
}

func (s *LocalStorageService) UploadFromFile(ctx context.Context, key string, contentType string, srcPath string) error {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	dir := filepath.Dir(key)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(key, data, 0644)
}
