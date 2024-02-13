package imageHandler

import (
	"image/jpeg"
	"io"
	"log"

	"github.com/adrium/goheif"
)

func ResizeImage(file io.Reader, resizedFile io.Writer) {

	img, err := goheif.Decode(file)

	if err != nil {
		log.Fatal(err)
	}

	if err := jpeg.Encode(resizedFile, img, nil); err != nil {
		log.Fatal(err)
	}

}
