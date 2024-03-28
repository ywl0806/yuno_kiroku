// Package imageHandler provides functions for handling image files.

package imageHandler

import (
	"bytes"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"log"
	"log/slog"
	"strings"

	"io"

	"github.com/labstack/echo/v4"
	"github.com/nfnt/resize"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/ywl0806/yuno_kiroku/api/utils"
)

type Imagehandler struct {
	OriginalFile  io.Reader
	ResizedFile   io.Writer
	OriginalImage image.Image
	ResizedImage  image.Image

	Ext  string
	Exif *exif.Exif
}

func NewImageHandler(originalFile io.Reader, resizedFile io.Writer, ext string) *Imagehandler {
	smallExt := strings.ToLower(ext)
	handler := &Imagehandler{
		OriginalFile: originalFile,
		Ext:          smallExt,
		ResizedFile:  resizedFile,
	}

	handler.decodeImage()
	return handler
}

// ResizeImage resizes the given image file to a suitable size for mobile devices.
// It takes the original image file as input and writes the resized image to the specified output file.
// The ext parameter specifies the file extension of the image.
// It returns the exif data of the original image and an error if any error occurs during the resizing process.
func (ih *Imagehandler) ResizeImage(maxWidth, maxHeight uint) (err error) {

	// Resize the image to a suitable size for mobile devices
	// max width 1500px, max height 1500px
	ih.ResizedImage = resize.Thumbnail(maxWidth, maxHeight, ih.OriginalImage, resize.Lanczos3)

	// Encode the image to jpeg format
	if err := jpeg.Encode(ih.ResizedFile, ih.ResizedImage, nil); err != nil {
		log.Println("image encode error: ", err)
		return echo.NewHTTPError(500, "image encode error")
	}

	return err
}

// decode image
func (ih *Imagehandler) decodeImage() (err error) {
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
func (ih *Imagehandler) decodeHeicImage() (err error) {
	file := new(bytes.Buffer)
	file, ih.OriginalFile, _ = utils.CopyReader(ih.OriginalFile)

	var exifsBytes []byte
	var exifsBuffer *bytes.Buffer

	ih.OriginalImage, exifsBytes, err = handleHeic(file)

	if err != nil {
		slog.Error("heic decode error: ", err)
		return echo.NewHTTPError(500, "heic decode error")
	}
	exifsBuffer = new(bytes.Buffer)
	NewWriterExif(exifsBuffer, exifsBytes)

	// Write the exif data to the resized file
	ih.ResizedFile, _ = NewWriterExif(ih.ResizedFile, exifsBytes)

	ih.Exif, err = exif.Decode(exifsBuffer)
	if err != nil {
		slog.Warn("exif decode error: ", err)
		err = nil
		ih.Exif = &exif.Exif{}
	}

	return
}

// decode "jpeg", "jpg", "png", "gif" image
func (ih *Imagehandler) decodeNomalImage() (err error) {
	file := new(bytes.Buffer)
	file, ih.OriginalFile, _ = utils.CopyReader(ih.OriginalFile)

	exifFile := new(bytes.Buffer)
	exifFile, file, _ = utils.CopyReader(file)

	ih.OriginalImage, _, err = image.Decode(file)

	if err != nil {
		slog.Error("image decode error: ", err)
		return echo.NewHTTPError(500, "image decode error")
	}

	ih.Exif, err = exif.Decode(exifFile)
	if err != nil {
		slog.Warn("exif decode error: ", err)
		err = nil
		ih.Exif = &exif.Exif{}
	}

	return
}
