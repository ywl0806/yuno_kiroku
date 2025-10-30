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

	"io"

	"github.com/disintegration/imaging"
	"github.com/labstack/echo/v4"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/ywl0806/yuno_kiroku/internal/api/utils"
)

type ImageHelper struct {
	OriginalFile  io.Reader
	ResizedFile   io.Writer
	OriginalImage image.Image
	ResizedImage  image.Image

	Ext  string
	Exif *exif.Exif
}

func NewImageHandler(originalFile io.Reader, resizedFile io.Writer, ext string) *ImageHelper {
	smallExt := strings.ToLower(ext)
	handler := &ImageHelper{
		OriginalFile: originalFile,
		Ext:          smallExt,
		ResizedFile:  resizedFile,
	}

	handler.decodeImage()
	return handler
}

// @description 이미지 리사이즈
// @param maxWidth uint
// @param maxHeight uint
// @return err error
func (ih *ImageHelper) ResizeImage(maxWidth, maxHeight uint) (err error) {

	// 이미지를 적절한 크기로 리사이즈
	// max width 1500px, max height 1500px
	ih.ResizedImage = imaging.Resize(ih.OriginalImage, int(maxWidth), int(maxHeight), imaging.Lanczos)

	// 이미지를 jpeg 포맷으로 인코딩
	if err := jpeg.Encode(ih.ResizedFile, ih.ResizedImage, nil); err != nil {
		log.Println("image encode error: ", err)
		return echo.NewHTTPError(500, "image encode error")
	}

	return err
}

// @description 이미지 디코딩
// @param err error
// @return err error
func (ih *ImageHelper) decodeImage() (err error) {
	switch ih.Ext {
	case "jpeg", "jpg", "png", "gif":
		err = ih.decodeNomalImage()
	case "heic", "heif":
		err = ih.decodeHeicImage()
	default:
		err = echo.NewHTTPError(400, "unsupported file type")
	}

	return
}

// decode heic, heif image
func (ih *ImageHelper) decodeHeicImage() (err error) {
	file := new(bytes.Buffer)
	file, ih.OriginalFile, _ = utils.CopyReader(ih.OriginalFile)

	var exifsBytes []byte
	var exifsBuffer *bytes.Buffer

	// heic 이미지 디코딩
	ih.OriginalImage, exifsBytes, err = handleHeic(file)

	if err != nil {
		log.Println("heic decode error: ", err)
		return echo.NewHTTPError(500, "heic decode error")
	}
	exifsBuffer = new(bytes.Buffer)
	NewWriterExif(exifsBuffer, exifsBytes)

	// exif 데이터를 리사이즈된 파일에 쓰기
	ih.ResizedFile, _ = NewWriterExif(ih.ResizedFile, exifsBytes)

	ih.Exif, err = exif.Decode(exifsBuffer)
	if err != nil {
		log.Println("exif decode error: ", err)
		err = nil
		ih.Exif = &exif.Exif{}
	}

	return
}

// @description "jpeg", "jpg", "png", "gif" 이미지를 디코딩
// @param err error
// @return err error
func (ih *ImageHelper) decodeNomalImage() (err error) {
	file := new(bytes.Buffer)
	file, ih.OriginalFile, _ = utils.CopyReader(ih.OriginalFile)

	exifFile := new(bytes.Buffer)
	exifFile, file, _ = utils.CopyReader(file)

	ih.OriginalImage, _, err = image.Decode(file)

	if err != nil {
		log.Println("image decode error: ", err)
		return echo.NewHTTPError(500, "image decode error")
	}

	ih.Exif, err = exif.Decode(exifFile)
	if err != nil {
		log.Println("exif decode error: ", err)
		err = nil
		ih.Exif = &exif.Exif{}
	}

	return
}
