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

// ResizeImage resizes the given image file to a suitable size for mobile devices.
// It takes the original image file as input and writes the resized image to the specified output file.
// The ext parameter specifies the file extension of the image.
// It returns the exif data of the original image and an error if any error occurs during the resizing process.
func ResizeImage(originalFile io.Reader, resizedFile io.Writer, ext string) (exifData *exif.Exif, err error) {
	var img image.Image

	smallExt := strings.ToLower(ext)
	exifFile, file, _ := utils.CopyReader(originalFile)

	// If the image is in jpeg, jpg, png, or gif format, use the standard library to decode it
	switch smallExt {
	case "jpeg", "jpg", "png", "gif":
		img, _, err = image.Decode(file)

		if err != nil {
			slog.Error("image decode error: ", err)
			return nil, echo.NewHTTPError(500, "image decode error")
		}

		exifData, err = exif.Decode(exifFile)
		if err != nil {
			slog.Warn("exif decode error: ", err)
			err = nil
			exifData = &exif.Exif{}
		}

	case "heic", "heif":
		var exifsBytes []byte
		var exifsBuffer *bytes.Buffer

		img, exifsBytes, err = handleHeic(file)

		if err != nil {
			slog.Error("heic decode error: ", err)
			return nil, echo.NewHTTPError(500, "heic decode error")
		}
		exifsBuffer = new(bytes.Buffer)
		NewWriterExif(exifsBuffer, exifsBytes)
		// Write the exif data to the resized file
		resizedFile, _ = NewWriterExif(resizedFile, exifsBytes)

		exifData, err = exif.Decode(exifsBuffer)
		if err != nil {
			slog.Warn("exif decode error: ", err)
			err = nil
			exifData = &exif.Exif{}
		}

	default:
		log.Println("unsupported file type")
		return nil, echo.NewHTTPError(400, "unsupported file type")
	}

	// Resize the image to a suitable size for mobile devices
	// max width 1500px, max height 1500px
	resizedImg := resize.Thumbnail(1500, 1500, img, resize.Lanczos3)

	// Encode the image to jpeg format
	if err := jpeg.Encode(resizedFile, resizedImg, nil); err != nil {
		log.Println("image encode error: ", err)
		return nil, echo.NewHTTPError(500, "image encode error")
	}

	return exifData, err
}
