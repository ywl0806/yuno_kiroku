package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ywl0806/yuno_kiroku/internal/consts"
	imagepkg "github.com/ywl0806/yuno_kiroku/pkg/image"
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
type ImageUploader struct {
	storage storage.StorageService
	cropper imagepkg.Cropper
}

func NewImageUploader(storage storage.StorageService, cropper imagepkg.Cropper) *ImageUploader {
	return &ImageUploader{storage: storage, cropper: cropper}
}

// 원본 이미지용 S3 키를 미리 생성 (Presigned URL 발급 시 사용)
// 반환 형식: {familyId}/{albumId}/{date}/original/{uuid}.{ext}
func (u *ImageUploader) BuildOriginalKey(familyId, albumId int32, fileName string, date time.Time) string {
	ext := strings.ToLower(strings.TrimPrefix(fileName[strings.LastIndex(fileName, "."):], "."))
	uploadPath := strconv.Itoa(int(familyId)) + "/" + strconv.Itoa(int(albumId))
	folderName := uploadPath + "/" + date.Format("2006-01-02")
	return folderName + "/" + consts.ORIGINAL_STORAGE_PREFIX + "/" + uuid.New().String() + "." + ext
}

// Presigned PUT URL을 생성
func (u *ImageUploader) GeneratePresignedPutURL(key string, contentType string, expiresIn time.Duration) (string, error) {
	return u.storage.GeneratePresignedPutURL(key, contentType, expiresIn)
}

// 가족/앨범에 대한 스토리지 경로 prefix를 반환
func (u *ImageUploader) BuildUploadPath(familyId, albumId int32) string {
	path := strconv.Itoa(int(familyId))
	if albumId != 0 {
		path += "/" + strconv.Itoa(int(albumId))
	}
	return path
}

// 원본 이미지 bytes를 스토리지에 저장
func (u *ImageUploader) UploadOriginal(data []byte, uploadPath string, meta imagepkg.Meta) (string, error) {
	folderName := uploadPath + "/" + meta.TakenAt.Format("2006-01-02")
	filename := uuid.New().String()
	storageKey, err := u.storage.SaveFile(data, folderName, consts.ORIGINAL_STORAGE_PREFIX+"/"+filename+"."+meta.Ext)
	if err != nil {
		return "", fmt.Errorf("원본 이미지 업로드 실패: %w", err)
	}
	return storageKey, nil
}

// 리사이즈된 이미지를 지정한 prefix 아래 스토리지에 저장
func (u *ImageUploader) UploadResized(img imagepkg.ResizedImage, uploadPath, prefix string, takenAt time.Time) (UploadResult, error) {
	folderName := uploadPath + "/" + takenAt.Format("2006-01-02")
	filename := uuid.New().String()
	storageKey, err := u.storage.SaveFile(img.Data, folderName, prefix+"/"+filename+img.Ext)
	if err != nil {
		return UploadResult{}, fmt.Errorf("%s 이미지 업로드 실패: %w", prefix, err)
	}
	return UploadResult{
		StorageKey: storageKey,
		Width:      img.Width,
		Height:     img.Height,
		FileSize:   int64(len(img.Data)),
	}, nil
}

// 얼굴 bbox에 패딩을 적용한 크롭 영역을 계산하고,
// cropper로 크롭·리사이즈한 뒤 스토리지에 업로드합니다.
func (u *ImageUploader) CropFaceAndUpload(
	viewData []byte, imgWidth, imgHeight int,
	top, right, bottom, left int32, padding float64,
	identityID, mediaItemID int32,
) (string, error) {
	// bbox 패딩 계산
	faceW := int(right - left)
	faceH := int(bottom - top)

	padX := int(float64(faceW) * padding)
	padY := int(float64(faceH) * padding)

	if faceW > faceH {
		padY += int(float64(faceW-faceH)*padding) + (faceW-faceH)/2
	} else {
		padX += int(float64(faceH-faceW)*padding) + (faceH-faceW)/2
	}

	cropLeft := clampMin(int(left)-padX, 0)
	cropTop := clampMin(int(top)-padY, 0)
	cropRight := clampMax(int(right)+padX, imgWidth)
	cropBottom := clampMax(int(bottom)+padY, imgHeight)

	cropW := cropRight - cropLeft
	cropH := cropBottom - cropTop

	croppedBytes, err := u.cropper.CropAndResize(viewData, cropLeft, cropTop, cropW, cropH, 512)
	if err != nil {
		return "", fmt.Errorf("얼굴 크롭 실패: %w", err)
	}
	defer func() { croppedBytes = nil }()

	storageKey, err := u.storage.SaveFile(
		croppedBytes,
		fmt.Sprintf("identities/%d", identityID),
		fmt.Sprintf("face_%d.webp", mediaItemID),
	)
	if err != nil {
		return "", fmt.Errorf("얼굴 크롭 이미지 업로드 실패: %w", err)
	}

	return storageKey, nil
}

func clampMin(v, min int) int {
	if v < min {
		return min
	}
	return v
}

func clampMax(v, max int) int {
	if v > max {
		return max
	}
	return v
}
