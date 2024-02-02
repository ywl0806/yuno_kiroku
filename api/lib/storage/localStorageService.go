package storage

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
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

func (s *LocalStorageService) SaveFile(file *multipart.FileHeader, filePath string) (string, error) {
	src, err := file.Open()

	if err != nil {
		return "", err
	}

	defer src.Close()

	path := filepath.Join("uploads", s.rootDir, filePath, file.Filename)
	dst, err := os.Create(path)

	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}
	return path, err
}