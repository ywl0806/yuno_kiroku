// Package imageHandler provides functions for handling image files.

package imageHandler

import (
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"log"
	"strings"

	"io"

	"github.com/adrium/goheif"
	"github.com/labstack/echo/v4"
	"github.com/nfnt/resize"
)

// ResizeImage resizes the given image file to a suitable size for mobile devices.
// It takes the original image file as input and writes the resized image to the specified output file.
// The ext parameter specifies the file extension of the image.
// It returns an error if any error occurs during the resizing process.
func ResizeImage(file io.Reader, resizedFile io.Writer, ext string) (err error) {
	var img image.Image
	smallExt := strings.ToLower(ext)

	// If the image is in jpeg, jpg, png, or gif format, use the standard library to decode it

	switch smallExt {
	case "jpeg", "", "png", "gif":
		img, _, err = image.Decode(file)
		if err != nil {
			log.Println("image decode error: ", err)
			return echo.NewHTTPError(500, "image decode error")
		}
	case "heic", "heif":
		img, err = handleHeic(file)
		if err != nil {
			log.Println("heic decode error: ", err)
			return echo.NewHTTPError(500, "heic decode error")
		}
		// The heic image has already been handled
	default:
		log.Println("unsupported file type")
		return echo.NewHTTPError(400, "unsupported file type")
	}

	// Resize the image to a suitable size for mobile devices
	// max width 1500px, max height 1500px
	resizedImg := resize.Thumbnail(1500, 1500, img, resize.Lanczos3)

	// Encode the image to jpeg format
	if err := jpeg.Encode(resizedFile, resizedImg, nil); err != nil {
		log.Println("image encode error: ", err)
		return echo.NewHTTPError(500, "image encode error")
	}

	return nil
}

// handleHeic decodes the given heic image file using goheif library.
// It returns the decoded image and an error if any error occurs during the decoding process.
func handleHeic(file io.Reader) (image.Image, error) {
	img, err := goheif.Decode(file)
	if err != nil {
		log.Println("heic decode error: ", err)
		return nil, err
	}
	return img, nil
}
