// Package imageHandler provides functions for handling image files.

package imageHandler

import (
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"

	"io"

	"github.com/adrium/goheif"
	"github.com/nfnt/resize"
)

// ResizeImage resizes the given image file to a suitable size for mobile devices.
// It takes the original image file as input and writes the resized image to the specified output file.
// The ext parameter specifies the file extension of the image.
// It returns an error if any error occurs during the resizing process.
func ResizeImage(file io.Reader, resizedFile io.Writer, ext string) (err error) {
	var img image.Image

	// If the image is in heic format, use goheif to decode it
	if ext == "heic" || ext == "heif" {
		img, err = handleHeic(file)
		if err != nil {
			return err
		}
	}

	// If the image is in jpeg, jpg, png, or gif format, use the standard library to decode it
	if ext == "jpeg" || ext == "jpg" || ext == "png" || ext == "gif" {
		img, _, err = image.Decode(file)
		if err != nil {
			return err
		}
	}

	// Resize the image to a suitable size for mobile devices
	// max width 1500px, max height 1500px
	resizedImg := resize.Thumbnail(1500, 1500, img, resize.Lanczos3)

	// Encode the image to jpeg format
	if err := jpeg.Encode(resizedFile, resizedImg, nil); err != nil {
		return err
	}

	return nil
}

// handleHeic decodes the given heic image file using goheif library.
// It returns the decoded image and an error if any error occurs during the decoding process.
func handleHeic(file io.Reader) (image.Image, error) {
	img, err := goheif.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("heic decode error: %v", err)
	}
	return img, nil
}
