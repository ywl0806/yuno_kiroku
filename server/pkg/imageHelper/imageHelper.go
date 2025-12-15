// Package imageHandler provides functions for handling image files.

package imageHandler

import (
	"bytes"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"log"
	"strings"

	"io"

	"github.com/disintegration/imaging"
	"github.com/labstack/echo/v4"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/ywl0806/yuno_kiroku/internal/api/utils"
)

const (
	MaxWidth = 1500
)

type ImageHelper struct {
	OriginalFile  io.ReadSeeker
	ResizedFile   *bytes.Buffer // bytes.Buffer는 io.ReadWriter를 구현 (Reader + Writer 둘 다 가능)
	OriginalImage image.Image
	ResizedImage  image.Image

	Ext  string
	Exif *exif.Exif
}

func NewImageHandler(originalFile io.ReadSeeker, ext string) *ImageHelper {
	smallExt := strings.ToLower(ext)
	resizedFile := bytes.NewBuffer([]byte{})
	handler := &ImageHelper{
		OriginalFile: originalFile,
		Ext:          smallExt,
		ResizedFile:  resizedFile,
	}

	handler.decodeImage()
	handler.resizeImage()
	return handler
}

// 이미지 리사이즈
func (ih *ImageHelper) resizeImage() (err error) {

	// 이미지를 적절한 크기로 리사이즈
	// max width 1500px, max height 1500px
	ih.ResizedImage = imaging.Resize(ih.OriginalImage, MaxWidth, 0, imaging.Lanczos)

	// 이미지를 jpeg 포맷으로 인코딩 (exifWriter에 쓰면 ResizedFile 버퍼에 저장됨)
	if err := jpeg.Encode(ih.ResizedFile, ih.ResizedImage, nil); err != nil {
		log.Println("image encode error: ", err)
		return err
	}

	_, err = NewWriterExif(ih.ResizedFile, ih.Exif.Raw)
	if err != nil {
		log.Println("exif write error: ", err)
		return err
	}

	return nil
}

// 이미지 디코딩
func (ih *ImageHelper) decodeImage() (err error) {
	var img image.Image
	var exifData *exif.Exif

	switch ih.Ext {
	case "jpeg", "jpg", "png", "gif":
		img, exifData, err = ih.decodeNomalImage()
		if err != nil {
			return err
		}
	case "heic", "heif":
		img, exifData, err = ih.decodeHeicImage()
		if err != nil {
			return err
		}
	default:
		err = echo.NewHTTPError(400, "unsupported file type")
	}
	ih.OriginalImage = img
	ih.Exif = exifData

	return nil
}

// HEIC 이미지 디코딩
func (ih *ImageHelper) decodeHeicImage() (image.Image, *exif.Exif, error) {
	// 한 번만 읽어서 버퍼 생성 (handleHeic 내부에서 TeeReader 사용)
	fileBuf, _, err := utils.CopyReader(ih.OriginalFile)
	if err != nil {
		log.Println("copy reader error: ", err)
		return nil, nil, err
	}
	// heic 이미지 디코딩 (handleHeic가 EXIF도 함께 추출)
	originalImage, exifBuf, err := handleHeic(fileBuf)
	if err != nil {
		log.Println("heic decode error: ", err)
		return nil, nil, err
	}
	// exif 디코딩
	exifData, err := exif.Decode(bytes.NewReader(exifBuf))
	if err != nil {
		log.Println("exif decode error: ", err)
		// EXIF 디코딩 실패는 치명적이지 않으므로 nil로 처리
		return originalImage, nil, nil
	}
	return originalImage, exifData, nil
}

// "jpeg", "jpg", "png", "gif" 이미지를 디코딩
func (ih *ImageHelper) decodeNomalImage() (image.Image, *exif.Exif, error) {
	// 한 번만 읽어서 이미지용과 EXIF용 두 개의 버퍼 생성
	imageBuf, exifBuf, err := utils.CopyReader(ih.OriginalFile)
	ih.OriginalFile.Seek(0, io.SeekStart)

	if err != nil {
		log.Println("copy reader error: ", err)
		return nil, nil, err
	}
	// 이미지 디코딩
	originalImage, _, err := image.Decode(imageBuf)
	if err != nil {
		log.Println("image decode error: ", err)
		return nil, nil, err
	}
	// EXIF 디코딩 (실패해도 이미지는 사용 가능)
	exifData, err := exif.Decode(bytes.NewReader(exifBuf.Bytes()))
	if err != nil {
		log.Println("exif decode error: ", err)
		// EXIF 디코딩 실패는 치명적이지 않으므로 nil로 처리
		return originalImage, nil, nil
	}
	return originalImage, exifData, nil
}
func (ih *ImageHelper) GetDominantColor(img image.Image) color.RGBA {
	var r, g, b, count float64
	rect := img.Bounds()
	for i := 0; i < rect.Max.Y; i++ {
		for j := 0; j < rect.Max.X; j++ {
			c := color.RGBAModel.Convert(img.At(j, i))
			r += float64(c.(color.RGBA).R)
			g += float64(c.(color.RGBA).G)
			b += float64(c.(color.RGBA).B)
			count++
		}
	}
	return color.RGBA{
		R: uint8(r / count),
		G: uint8(g / count),
		B: uint8(b / count),
	}
}
