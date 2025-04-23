package imageHandler

import (
	"bytes"
	"image"
	"io"
	"log"

	"github.com/adrium/goheif"
)

// handleHeic decodes the given heic image file using goheif library.
// It returns the decoded image and an error if any error occurs during the decoding process.
func handleHeic(file io.Reader) (image.Image, []byte, error) {

	exifFile := new(bytes.Buffer)
	tee := io.TeeReader(file, exifFile)

	img, err := goheif.Decode(tee)
	if err != nil {
		log.Println("heic img decode error: ", err)
		return nil, nil, err
	}

	buf := new(bytes.Buffer)
	io.Copy(buf, exifFile)
	reader := bytes.NewReader(buf.Bytes())

	heifExif, err := goheif.ExtractExif(reader)

	if err != nil {
		log.Println("heic exif error: ", err)
		return nil, nil, err
	}

	return img, heifExif, nil
}
