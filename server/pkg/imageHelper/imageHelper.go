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

type ResizedFile struct {
	File   []byte
	Width  int
	Height int
	// 리사이즈된 이미지 확장자 (ex: .jpeg, .png, .webp ...)
	Ext string
}
type ImageHelper struct {
	OriginalFile   []byte
	OriginalWidth  int
	OriginalHeight int

	resizedFiles map[int]ResizedFile

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
		Ext:          smallExt,
		resizedFiles: make(map[int]ResizedFile),
	}

	originalFileBytes, err := io.ReadAll(originalFile)
	if err != nil {
		return nil, err
	}
	handler.OriginalFile = originalFileBytes

	// 원본 이미지 크기 가져오기
	originalWidth, originalHeight, err := handler.getOriginalImageSize()
	if err != nil {
		return nil, err
	}
	handler.OriginalWidth = originalWidth
	handler.OriginalHeight = originalHeight

	return handler, nil
}

func (ih *ImageHelper) resize(maxLength int, ext string) (ResizedFile, error) {

	newWidth, newHeight := ih.calculateResizedImageSize(ih.OriginalWidth, ih.OriginalHeight, maxLength)

	if newWidth == ih.OriginalWidth && newHeight == ih.OriginalHeight {
		return ResizedFile{
			File:   ih.OriginalFile,
			Width:  ih.OriginalWidth,
			Height: ih.OriginalHeight,
			Ext:    ih.Ext,
		}, nil
	}

	img, err := vips.NewImageFromBuffer(ih.OriginalFile)

	if err != nil {
		return ResizedFile{}, err
	}
	defer img.Close()

	// EXIF orientation 정보에 따라 자동 회전
	err = img.AutoRotate()
	if err != nil {
		return ResizedFile{}, err
	}

	actualWidth := img.Width()

	// 리사이즈 (비율 계산) - 회전 후 실제 크기 기준
	scale := float64(newWidth) / float64(actualWidth)
	err = img.Resize(scale, vips.KernelLanczos3)
	if err != nil {
		return ResizedFile{}, err
	}

	exportParams := vips.NewWebpExportParams()
	exportParams.Quality = 90
	exportParams.StripMetadata = true

	// WEBP로 내보내기 (EXIF 포함)
	resizedBytes, metadata, err := img.ExportWebp(exportParams)
	if err != nil {
		return ResizedFile{}, err
	}

	return ResizedFile{
		File:   resizedBytes,
		Width:  metadata.Width,
		Height: metadata.Height,
		Ext:    metadata.Format.FileExt(),
	}, nil
}

// 리사이즈된 이미지 가져오기(없으면 리사이즈 후 저장)
func (ih *ImageHelper) GetResizedFile(maxLength int) (*ResizedFile, error) {
	if resizedFile, ok := ih.resizedFiles[maxLength]; ok {
		return &resizedFile, nil
	}
	resizedFile, err := ih.resize(maxLength, ih.Ext)
	if err != nil {
		return nil, err
	}
	ih.resizedFiles[maxLength] = resizedFile
	return &resizedFile, nil
}

func (ih *ImageHelper) getOriginalImageSize() (int, int, error) {
	// libvips로 이미지 크기 가져오기
	img, err := vips.NewImageFromBuffer(ih.OriginalFile)
	if err != nil {
		log.Println("vips load image for size error: ", err)
		return 0, 0, err
	}
	defer img.Close()

	// EXIF orientation 정보에 따라 자동 회전 후 크기 가져오기
	err = img.AutoRotate()
	if err != nil {
		log.Println("vips auto rotate error (size check): ", err)
		// 회전 실패해도 원본 크기 사용
	}

	// 회전 후 실제 이미지 크기 저장
	originalWidth := img.Width()
	originalHeight := img.Height()
	return originalWidth, originalHeight, nil
}

func (ih *ImageHelper) GetOriginalImageSize() (int, int) {
	return ih.OriginalWidth, ih.OriginalHeight
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
func (ih *ImageHelper) calculateResizedImageSize(originalWidth int, originalHeight int, maxLength int) (int, int) {

	if maxLength > originalWidth && maxLength > originalHeight {
		return originalWidth, originalHeight
	}

	newWidth := maxLength
	newHeight := maxLength

	if originalWidth > originalHeight {
		newHeight = int(float64(originalHeight) * float64(maxLength) / float64(originalWidth))
	} else {
		newWidth = int(float64(originalWidth) * float64(maxLength) / float64(originalHeight))
	}

	return newWidth, newHeight
}
