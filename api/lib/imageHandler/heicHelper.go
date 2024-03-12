package imageHandler

import (
	"bytes"
	"image"
	"io"
	"log"

	"github.com/adrium/goheif"
)

// Skip Writer for exif writing
type writerSkipper struct {
	w           io.Writer
	bytesToSkip int
}

func (w *writerSkipper) Write(data []byte) (int, error) {
	if w.bytesToSkip <= 0 {
		return w.w.Write(data)
	}

	if dataLen := len(data); dataLen < w.bytesToSkip {
		w.bytesToSkip -= dataLen
		return dataLen, nil
	}

	if n, err := w.w.Write(data[w.bytesToSkip:]); err == nil {
		n += w.bytesToSkip
		w.bytesToSkip = 0
		return n, nil
	} else {
		return n, err
	}
}

func newWriterExif(w io.Writer, exif []byte) (io.Writer, error) {
	writer := &writerSkipper{w, 2}
	soi := []byte{0xff, 0xd8}
	if _, err := w.Write(soi); err != nil {
		return nil, err
	}

	if exif != nil {
		app1Marker := 0xe1
		markerlen := 2 + len(exif)
		marker := []byte{0xff, uint8(app1Marker), uint8(markerlen >> 8), uint8(markerlen & 0xff)}
		if _, err := w.Write(marker); err != nil {
			return nil, err
		}

		if _, err := w.Write(exif); err != nil {
			return nil, err
		}
	}

	return writer, nil
}

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

	exif, err := goheif.ExtractExif(reader)

	if err != nil {
		log.Println("heic exif error: ", err)
		return nil, nil, err
	}

	return img, exif, nil
}
