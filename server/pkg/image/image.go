// Package image는 libvips를 사용한 이미지 처리 함수를 제공합니다.
package image

import (
	"sync"
	"time"

	"github.com/davidbyttow/govips/v2/vips"
)

var vipsOnce sync.Once

// libvips 초기화
func initVips() {
	vipsOnce.Do(func() {
		vips.Startup(nil)
	})
}

// 원본 파일에서 추출한 이미지 메타데이터
type Meta struct {
	Width    int
	Height   int
	FileSize int64
	TakenAt  time.Time
	Lat      *float64
	Lon      *float64
	Ext      string
}

// 리사이즈 결과를 담는 구조체
type ResizedImage struct {
	Data   []byte
	Width  int
	Height int
	Ext    string // 예: ".webp"
}

// 이미지 영역을 크롭·리사이즈하는 인터페이스
// DI/테스트를 위해 분리되어 있음
type Cropper interface {
	CropAndResize(data []byte, left, top, width, height, targetSize int) ([]byte, error)
}

// 이미지 bytes에서 메타데이터(크기, EXIF, GPS)를 추출
func Parse(data []byte, ext string) (Meta, error) {
	initVips()

	img, err := vips.NewImageFromBuffer(data)
	if err != nil {
		return Meta{}, err
	}
	defer img.Close()

	// EXIF orientation에 따라 실제 표시 방향으로 회전 후 크기 측정 (AutoRotate)
	_ = img.AutoRotate()

	takenAt, lat, lon := extractExif(data)

	return Meta{
		Width:    img.Width(),
		Height:   img.Height(),
		FileSize: int64(len(data)),
		TakenAt:  takenAt,
		Lat:      lat,
		Lon:      lon,
		Ext:      ext,
	}, nil
}

// 긴 변이 maxLength가 되도록 비율을 유지하며 리사이즈
// 결과는 WebP 형식으로 반환
func Resize(data []byte, maxLength int) (ResizedImage, error) {
	initVips()

	img, err := vips.NewImageFromBuffer(data)
	if err != nil {
		return ResizedImage{}, err
	}
	defer img.Close()

	if err = img.AutoRotate(); err != nil {
		return ResizedImage{}, err
	}

	origW, origH := img.Width(), img.Height()
	newW, newH := calcResizedDimensions(origW, origH, maxLength)

	// 원본 크기 이하면 리사이즈 불필요
	if newW == origW && newH == origH {
		exportParams := vips.NewWebpExportParams()
		exportParams.Quality = 90
		exportParams.StripMetadata = true
		resizedBytes, metadata, err := img.ExportWebp(exportParams)
		if err != nil {
			return ResizedImage{}, err
		}
		return ResizedImage{
			Data:   resizedBytes,
			Width:  metadata.Width,
			Height: metadata.Height,
			Ext:    metadata.Format.FileExt(),
		}, nil
	}

	scale := float64(newW) / float64(origW)
	if err = img.Resize(scale, vips.KernelLanczos3); err != nil {
		return ResizedImage{}, err
	}

	exportParams := vips.NewWebpExportParams()
	exportParams.Quality = 90
	exportParams.StripMetadata = true

	resizedBytes, metadata, err := img.ExportWebp(exportParams)
	if err != nil {
		return ResizedImage{}, err
	}

	return ResizedImage{
		Data:   resizedBytes,
		Width:  metadata.Width,
		Height: metadata.Height,
		Ext:    metadata.Format.FileExt(),
	}, nil
}

// 지정한 영역을 크롭한 후 targetSize×targetSize로 리사이즈해 WebP로 반환
// 좌표 계산(얼굴 패딩 등)은 호출부에서 처리
func CropAndResize(data []byte, left, top, width, height, targetSize int) ([]byte, error) {
	initVips()

	img, err := vips.NewImageFromBuffer(data)
	if err != nil {
		return nil, err
	}
	defer img.Close()

	if err = img.ExtractArea(left, top, width, height); err != nil {
		return nil, err
	}

	scale := float64(targetSize) / float64(img.Width())
	if err = img.Resize(scale, vips.KernelLanczos3); err != nil {
		return nil, err
	}

	exportParams := vips.NewWebpExportParams()
	exportParams.Quality = 85
	exportParams.StripMetadata = true

	croppedBytes, _, err := img.ExportWebp(exportParams)
	if err != nil {
		return nil, err
	}

	return croppedBytes, nil
}

// CropAndResize 함수를 사용하는 기본 Cropper 구현체
type defaultCropper struct{}

func (c *defaultCropper) CropAndResize(data []byte, left, top, width, height, targetSize int) ([]byte, error) {
	return CropAndResize(data, left, top, width, height, targetSize)
}

// 기본 Cropper 구현체를 반환
func NewCropper() Cropper {
	return &defaultCropper{}
}

// 긴 변이 maxLength가 되도록 새 크기를 계산
func calcResizedDimensions(origW, origH, maxLength int) (int, int) {
	if maxLength >= origW && maxLength >= origH {
		return origW, origH
	}
	if origW > origH {
		return maxLength, int(float64(origH) * float64(maxLength) / float64(origW))
	}
	return int(float64(origW) * float64(maxLength) / float64(origH)), maxLength
}
