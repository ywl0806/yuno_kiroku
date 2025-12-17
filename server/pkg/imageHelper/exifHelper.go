package imageHandler

import (
	"io"
)

var (
	SOI = []byte{0xff, 0xd8}
	// APP1 Marker(jpeg exif 마커)
	APP1_MARKER = []byte{0xff, 0xe1}
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

func NewWriterExif(w io.Writer, exif []byte) (io.Writer, error) {
	writer := &writerSkipper{w, 2}

	// SOI(Start of Image) 마커 쓰기
	if _, err := w.Write(SOI); err != nil {
		return nil, err
	}

	// APP1 Marker(jpeg exif 마커) 쓰기
	if exif != nil {

		// APP1 Marker 길이 계산
		markerlen := 2 + len(exif)
		marker := []byte{APP1_MARKER[0], APP1_MARKER[1], uint8(markerlen >> 8), uint8(markerlen & 0xff)}
		if _, err := w.Write(marker); err != nil {
			return nil, err
		}

		// exif 데이터 쓰기
		if _, err := w.Write(exif); err != nil {
			return nil, err
		}
	}

	return writer, nil
}
