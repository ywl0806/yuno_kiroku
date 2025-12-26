// Package imageHandler provides functions for handling image files.

package imageHandler

import (
	"log"
	"strings"
	"sync"
	"time"

	"io"

	"github.com/davidbyttow/govips/v2/vips"
)

var (
	vipsOnce sync.Once
)

const (
	MaxLength = 2048
)

type ImageHelper struct {
	OriginalFile   []byte
	ResizedFile    []byte
	OriginalWidth  int
	OriginalHeight int
	ResizedWidth   int
	ResizedHeight  int

	ResizedOrientation int
	// 리사이즈된 이미지 확장자 (ex: .jpeg, .png, .webp ...)
	ResizedExt string

	Ext  string
	Exif map[string]string
}

func NewImageHandler(originalFile io.Reader, ext string) (*ImageHelper, error) {
	// libvips 초기화 (한 번만 실행)
	vipsOnce.Do(func() {
		vips.Startup(nil)
	})

	smallExt := strings.ToLower(ext)

	handler := &ImageHelper{
		Ext: smallExt,
	}

	originalFileBytes, err := io.ReadAll(originalFile)
	if err != nil {
		return nil, err
	}
	handler.OriginalFile = originalFileBytes

	// 원본 이미지 크기 가져오기
	err = handler.getOriginalImageSize()
	if err != nil {
		return nil, err
	}

	// 이미지 리사이즈
	err = handler.resizeImage()
	if err != nil {
		return nil, err
	}

	return handler, nil
}

// 이미지 리사이즈
func (ih *ImageHelper) resizeImage() (err error) {
	// 이미지를 적절한 크기로 리사이즈
	newWidth, newHeight := ih.calculateResizedImageSize(ih.OriginalWidth, ih.OriginalHeight)

	// 원본 크기와 같으면 리사이즈 불필요
	if newWidth == ih.OriginalWidth && newHeight == ih.OriginalHeight {
		ih.ResizedFile = ih.OriginalFile
		ih.ResizedWidth = ih.OriginalWidth
		ih.ResizedHeight = ih.OriginalHeight
		return nil
	}

	// libvips로 이미지 로드
	img, err := vips.NewImageFromBuffer(ih.OriginalFile)
	if err != nil {
		log.Println("vips load image error: ", err)
		return err
	}
	defer img.Close()

	// EXIF orientation 정보에 따라 자동 회전
	err = img.AutoRotate()
	if err != nil {
		log.Println("vips auto rotate error: ", err)
		// 회전 실패해도 계속 진행
	}

	ih.Exif = img.GetExif()

	// 회전 후 실제 이미지 크기 가져오기
	actualWidth := img.Width()

	// 리사이즈 (비율 계산) - 회전 후 실제 크기 기준
	scale := float64(newWidth) / float64(actualWidth)
	err = img.Resize(scale, vips.KernelLanczos3)
	if err != nil {
		log.Println("vips resize error: ", err)
		return err
	}

	// WEBP로 내보내기 (EXIF 포함)
	exportParams := vips.NewWebpExportParams()
	exportParams.Quality = 90
	exportParams.StripMetadata = true

	resizedBytes, metadata, err := img.ExportWebp(exportParams)
	if err != nil {
		log.Println("vips export error: ", err)
		return err
	}

	ih.ResizedFile = resizedBytes
	ih.ResizedExt = metadata.Format.FileExt()

	// 리사이즈 후 실제 크기 저장 (Width와 Height가 올바르게 설정됨)
	ih.ResizedWidth = metadata.Width
	ih.ResizedHeight = metadata.Height

	return nil
}

func (ih *ImageHelper) getOriginalImageSize() error {
	// libvips로 이미지 크기 가져오기
	img, err := vips.NewImageFromBuffer(ih.OriginalFile)
	if err != nil {
		log.Println("vips load image for size error: ", err)
		return err
	}
	defer img.Close()

	// EXIF orientation 정보에 따라 자동 회전 후 크기 가져오기
	err = img.AutoRotate()
	if err != nil {
		log.Println("vips auto rotate error (size check): ", err)
		// 회전 실패해도 원본 크기 사용
	}

	// 회전 후 실제 이미지 크기 저장
	ih.OriginalWidth = img.Width()
	ih.OriginalHeight = img.Height()
	return nil
}

func (ih *ImageHelper) GetOriginalImageSize() (int, int) {
	return ih.OriginalWidth, ih.OriginalHeight
}

func (ih *ImageHelper) GetResizedImageSize() (int, int) {
	return ih.ResizedWidth, ih.ResizedHeight
}

func (ih *ImageHelper) GetTakenAt() time.Time {
	takenAt, _ := ih.Exif["DateTime"]
	if takenAt == "" {
		return time.Now()
	}
	takenAtTime, err := time.Parse(time.RFC3339, takenAt)
	if err != nil {
		return time.Now()
	}
	return takenAtTime
}

// 긴변을 MaxLength로 고정하고 짧은변을 계산
func (ih *ImageHelper) calculateResizedImageSize(originalWidth int, originalHeight int) (int, int) {

	if MaxLength > originalWidth && MaxLength > originalHeight {
		return originalWidth, originalHeight
	}

	newWidth := MaxLength
	newHeight := MaxLength

	if originalWidth > originalHeight {
		newHeight = int(float64(originalHeight) * float64(MaxLength) / float64(originalWidth))
	} else {
		newWidth = int(float64(originalWidth) * float64(MaxLength) / float64(originalHeight))
	}

	return newWidth, newHeight
}
