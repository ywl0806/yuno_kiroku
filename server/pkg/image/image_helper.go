// Package imageHandler provides functions for handling image files.

package image

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"io"

	"github.com/davidbyttow/govips/v2/vips"
)

var (
	vipsOnce sync.Once
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
	Ext          string
	takenAt      time.Time
	latitude     *float64
	longitude    *float64
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

	// EXIF 정보 가져오기 (govips 사용 - HEIC 포함 지원)
	takenAt, lat, lon := extractExifFromVips(originalFileBytes)
	handler.takenAt = takenAt
	handler.latitude = lat
	handler.longitude = lon

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
	return ih.takenAt
}

func (ih *ImageHelper) GetLocation() (*float64, *float64) {
	return ih.latitude, ih.longitude
}

// CropFaceFromBytes: view 이미지 bytes에서 bbox 좌표로 얼굴 영역을 크롭 후 WebP bytes 반환
// padding: 얼굴 주변 여백 비율 (0.2 = 20%)
func CropFaceFromBytes(imageBytes []byte, top, right, bottom, left int, padding float64) ([]byte, error) {
	img, err := vips.NewImageFromBuffer(imageBytes)
	if err != nil {
		return nil, err
	}
	defer img.Close()

	imgWidth := img.Width()
	imgHeight := img.Height()

	faceWidth := right - left
	faceHeight := bottom - top

	padX := int(float64(faceWidth) * padding)
	padY := int(float64(faceHeight) * padding)

	if faceWidth > faceHeight {
		padY = padY + int((float64(faceWidth-faceHeight) * padding)) + (faceWidth-faceHeight)/2
	} else {
		padX = padX + int((float64(faceHeight-faceWidth) * padding)) + (faceHeight-faceWidth)/2
	}

	cropLeft := left - padX
	cropTop := top - padY
	cropRight := right + padX
	cropBottom := bottom + padY

	// 이미지 경계 클램핑
	if cropLeft < 0 {
		cropLeft = 0
	}
	if cropTop < 0 {
		cropTop = 0
	}
	if cropRight > imgWidth {
		cropRight = imgWidth
	}
	if cropBottom > imgHeight {
		cropBottom = imgHeight
	}

	cropWidth := cropRight - cropLeft
	cropHeight := cropBottom - cropTop

	err = img.ExtractArea(cropLeft, cropTop, cropWidth, cropHeight)
	if err != nil {
		return nil, err
	}

	// 512x512로 리사이즈
	const targetSize = 512
	scale := float64(targetSize) / float64(img.Width())
	err = img.Resize(scale, vips.KernelLanczos3)
	if err != nil {
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

// govips를 사용해 EXIF 메타데이터 추출 (HEIC 포함)
func extractExifFromVips(imageBytes []byte) (takenAt time.Time, lat *float64, lon *float64) {
	takenAt = time.Now()

	img, err := vips.NewImageFromBuffer(imageBytes)
	if err != nil {
		return
	}
	defer img.Close()

	exifData := img.GetExif()

	// 촬영 일시: DateTimeOriginal 우선, 없으면 DateTime
	for _, key := range []string{"exif-ifd2-DateTimeOriginal", "exif-ifd0-DateTime"} {
		if raw, ok := exifData[key]; ok {
			// 값 형식: "2025:11:26 13:33:59 (2025:11:26 13:33:59, ASCII, ...)"
			dateStr := extractExifValue(raw)
			if t, err := time.Parse("2006:01:02 15:04:05", dateStr); err == nil {
				takenAt = t
				break
			}
		}
	}

	// GPS 좌표
	latStr, latOk := exifData["exif-ifd3-GPSLatitude"]
	lonStr, lonOk := exifData["exif-ifd3-GPSLongitude"]
	if !latOk || !lonOk {
		return
	}

	latVal, err := parseGPSRational(extractExifValue(latStr))
	if err != nil {
		return
	}
	lonVal, err := parseGPSRational(extractExifValue(lonStr))
	if err != nil {
		return
	}

	if ref, ok := exifData["exif-ifd3-GPSLatitudeRef"]; ok {
		if extractExifValue(ref) == "S" {
			latVal = -latVal
		}
	}
	if ref, ok := exifData["exif-ifd3-GPSLongitudeRef"]; ok {
		if extractExifValue(ref) == "W" {
			lonVal = -lonVal
		}
	}

	lat = &latVal
	lon = &lonVal
	return
}

// govips EXIF 값에서 실제 값 추출 (예: "2025:11:26 13:33:59 (2025:11:26 13:33:59, ASCII, ...)" → "2025:11:26 13:33:59")
func extractExifValue(s string) string {
	s = strings.TrimRight(s, "\x00")
	if idx := strings.Index(s, " ("); idx >= 0 {
		return strings.TrimSpace(s[:idx])
	}
	return strings.TrimSpace(s)
}

// EXIF GPS rational 문자열 파싱 (예: "37/1 28/1 456/100" → 십진수 도)
func parseGPSRational(s string) (float64, error) {
	parts := strings.Fields(s)
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid GPS rational: %s", s)
	}

	divisors := []float64{1, 60, 3600}
	var result float64
	for i, part := range parts {
		sub := strings.SplitN(strings.TrimSpace(part), "/", 2)
		if len(sub) != 2 {
			return 0, fmt.Errorf("invalid rational: %s", part)
		}
		num, err1 := strconv.ParseFloat(sub[0], 64)
		den, err2 := strconv.ParseFloat(sub[1], 64)
		if err1 != nil || err2 != nil || den == 0 {
			return 0, fmt.Errorf("invalid rational components: %s", part)
		}
		result += (num / den) / divisors[i]
	}
	return result, nil
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
