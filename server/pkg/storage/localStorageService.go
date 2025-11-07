package storage

import (
	"io"
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

func (s *LocalStorageService) SaveFile(file io.Reader, filePath string, fileName string) (string, error) {

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

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return path, err
}
