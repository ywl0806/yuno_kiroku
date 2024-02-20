package imageHandler

import (
	"image"
	"image/jpeg"
	"io"
	"log"

	"github.com/adrium/goheif"
	"github.com/nfnt/resize"
)

func ResizeImage(file io.Reader, resizedFile io.Writer) {

	img, err := goheif.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	// Resize the image to a suitable size for mobile devices
	// max width 1500px, max height 1500px
	resizedImg := resize.Thumbnail(1500, 1500, img, resize.Lanczos3)

	if err := jpeg.Encode(resizedFile, resizedImg, nil); err != nil {
		log.Fatal(err)
	}

}

func handleHeic(file io.Reader) image.Image {
	img, err := goheif.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	return img
}
