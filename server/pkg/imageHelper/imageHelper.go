// Package imageHandler provides functions for handling image files.

package imageHandler

import (
	"bytes"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"log"
	"strings"
	"time"

	"io"

	"github.com/disintegration/imaging"
	"github.com/labstack/echo/v4"
	"github.com/rwcarlsen/goexif/exif"
)

const (
	MaxLength = 2048
)

type ImageHelper struct {
	OriginalFile  []byte
	ResizedFile   []byte
	OriginalImage image.Image
	ResizedImage  image.Image

	Ext       string
	Exif      *exif.Exif
	ExifBytes []byte // 원본 EXIF 바이트 데이터
}

func NewImageHandler(originalFile io.Reader, ext string) (*ImageHelper, error) {
	smallExt := strings.ToLower(ext)

	handler := &ImageHelper{
		Ext: smallExt,
	}

	originalImage, exifData, exifBytes, err := handler.decodeImage(originalFile)
	if err != nil {
		return nil, err
	}
	handler.OriginalImage = originalImage
	handler.Exif = exifData
	handler.ExifBytes = exifBytes
	handler.resizeImage()

	return handler, nil
}

// 이미지 리사이즈
func (ih *ImageHelper) resizeImage() (err error) {

	// 이미지를 적절한 크기로 리사이즈
	newWidth, newHeight := ih.calculateResizedImageSize(ih.GetOriginalImageSize())

	// max length 2048px
	resizedImage := imaging.Resize(ih.OriginalImage, newWidth, newHeight, imaging.Lanczos)

	ih.ResizedImage = resizedImage

	resizeFileBuffer := bytes.NewBuffer([]byte{})
	// TODO: heic이미지 이외에는 EXIF 데이터가 반영되지 않는 문제 수정
	// EXIF를 포함한 JPEG Writer 생성
	exifWriter, err := NewWriterExif(resizeFileBuffer, ih.ExifBytes)
	if err != nil {
		log.Println("exif writer error: ", err)
		return err
	}

	if err := jpeg.Encode(exifWriter, resizedImage, nil); err != nil {
		log.Println("image encode error: ", err)
		return err
	}

	ih.ResizedFile = resizeFileBuffer.Bytes()

	return nil
}

// 이미지 디코딩
func (ih *ImageHelper) decodeImage(fileReader io.Reader) (image.Image, *exif.Exif, []byte, error) {
	var img image.Image
	var exifData *exif.Exif
	var exifBytes []byte

	originalFileBytes, err := io.ReadAll(fileReader)
	if err != nil {
		log.Println("read file error: ", err)
		return nil, nil, nil, err
	}

	ih.OriginalFile = originalFileBytes

	switch ih.Ext {
	case "jpeg", "jpg", "png", "gif":
		img, exifData, exifBytes, err = ih.decodeNomalImage(originalFileBytes)
		if err != nil {
			return nil, nil, nil, err
		}
	case "heic", "heif":
		img, exifData, exifBytes, err = ih.decodeHeicImage(originalFileBytes)
		if err != nil {
			return nil, nil, nil, err
		}

	default:
		err = echo.NewHTTPError(400, "unsupported file type")
	}

	return img, exifData, exifBytes, nil
}

// HEIC 이미지 디코딩
func (ih *ImageHelper) decodeHeicImage(originalFileBytes []byte) (image.Image, *exif.Exif, []byte, error) {

	// heic 이미지 디코딩 (handleHeic가 EXIF도 함께 추출)
	originalImage, exifBytes, err := handleHeic(originalFileBytes)
	if err != nil {
		log.Println("heic decode error: ", err)
		return nil, nil, nil, err
	}

	// exif 디코딩
	exifData, err := exif.Decode(bytes.NewReader(exifBytes))
	if err != nil {
		log.Println("exif decode error: ", err)
		// EXIF 디코딩 실패는 치명적이지 않으므로 nil로 처리
		return originalImage, nil, nil, nil
	}
	return originalImage, exifData, exifBytes, nil
}

// "jpeg", "jpg", "png", "gif" 이미지를 디코딩
func (ih *ImageHelper) decodeNomalImage(originalFileBytes []byte) (image.Image, *exif.Exif, []byte, error) {
	// 이미지 디코딩
	originalImage, err := imaging.Decode(bytes.NewReader(originalFileBytes), imaging.AutoOrientation(true))
	if err != nil {
		log.Println("image decode error: ", err)
		return nil, nil, nil, err
	}
	// EXIF 디코딩 (실패해도 이미지는 사용 가능)
	exifData, err := exif.Decode(bytes.NewReader(originalFileBytes))
	if err != nil {
		log.Println("exif decode error: ", err)
		// EXIF 디코딩 실패는 치명적이지 않으므로 nil로 처리
		return originalImage, nil, nil, nil
	}
	return originalImage, exifData, exifData.Raw, nil
}

func (ih *ImageHelper) GetOriginalImageSize() (int, int) {
	return ih.OriginalImage.Bounds().Dx(), ih.OriginalImage.Bounds().Dy()
}

func (ih *ImageHelper) GetResizedImageSize() (int, int) {
	return ih.ResizedImage.Bounds().Dx(), ih.ResizedImage.Bounds().Dy()
}

func (ih *ImageHelper) GetPhotoCreatedAt() time.Time {
	photoCreatedAt, _ := ih.Exif.DateTime()
	if photoCreatedAt.IsZero() {
		return time.Now()
	}
	return photoCreatedAt
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
