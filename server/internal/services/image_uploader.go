package services

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ywl0806/yuno_kiroku/internal/consts"
	"github.com/ywl0806/yuno_kiroku/pkg/storage"
)

// 스토리지 업로드 결과를 담는 구조체
type UploadResult struct {
	StorageKey string
	Width      int
	Height     int
	FileSize   int64
}

// 스토리지 업로드와 얼굴 크롭 좌표 계산을 담당
type ImageUploader interface {
	BuildOriginalKey(familyId, mediaItemId int32, fileName string) string
	GeneratePresignedPutURL(ctx context.Context, key string, contentType string, expiresIn time.Duration) (string, error)
	UploadOriginal(ctx context.Context, data []byte, familyId, mediaItemId int32, ext string) (string, error)
	GetFile(ctx context.Context, key string) ([]byte, error)
	SaveFile(ctx context.Context, data []byte, folder, fileName string) (string, error)
}

type imageUploader struct {
	storage storage.StorageService
}

func NewImageUploader(storage storage.StorageService) ImageUploader {
	return &imageUploader{storage: storage}
}

// 원본 이미지용 S3 키를 생성 (Presigned URL 발급 시 사용)
// 반환 형식: original/{familyId}/{mediaItemId}.{ext}
func (u *imageUploader) BuildOriginalKey(familyId, mediaItemId int32, fileName string) string {
	ext := strings.ToLower(strings.TrimPrefix(fileName[strings.LastIndex(fileName, "."):], "."))
	return consts.ORIGINAL_STORAGE_PREFIX + "/" + strconv.Itoa(int(familyId)) + "/" + strconv.Itoa(int(mediaItemId)) + "." + ext
}

// Presigned PUT URL을 생성
func (u *imageUploader) GeneratePresignedPutURL(ctx context.Context, key string, contentType string, expiresIn time.Duration) (string, error) {
	return u.storage.GeneratePresignedPutURL(ctx, key, contentType, expiresIn)
}

// 원본 이미지 bytes를 스토리지에 저장
// 저장 경로: original/{familyId}/{mediaItemId}.{ext}
func (u *imageUploader) UploadOriginal(ctx context.Context, data []byte, familyId, mediaItemId int32, ext string) (string, error) {
	folder := consts.ORIGINAL_STORAGE_PREFIX + "/" + strconv.Itoa(int(familyId))
	fileName := strconv.Itoa(int(mediaItemId)) + "." + ext
	storageKey, err := u.storage.SaveFile(ctx, data, folder, fileName)
	if err != nil {
		return "", fmt.Errorf("원본 이미지 업로드 실패: %w", err)
	}
	return storageKey, nil
}

func (u *imageUploader) GetFile(ctx context.Context, key string) ([]byte, error) {
	return u.storage.GetFile(ctx, key)
}

func (u *imageUploader) SaveFile(ctx context.Context, data []byte, folder, fileName string) (string, error) {
	return u.storage.SaveFile(ctx, data, folder, fileName)
}
