package imageHandler

import (
	"bytes"
	"image"
	"io"
	"log"

	"github.com/adrium/goheif"
)

/*
*
@param file io.Reader
@return image.Image, []byte, error
@description HEIC이미지를 디코딩하고 EXIF 데이터를 추출합니다.
*/
func handleHeic(file io.Reader) (image.Image, []byte, error) {

	// heic 이미지의 exif 데이터 추출
	exifFile := new(bytes.Buffer)
	tee := io.TeeReader(file, exifFile)

	// heic 이미지 디코딩
	img, err := goheif.Decode(tee)
	if err != nil {
		log.Println("heic img decode error: ", err)
		return nil, nil, err
	}

	// 리더 생성
	buf := new(bytes.Buffer)
	io.Copy(buf, exifFile)
	reader := bytes.NewReader(buf.Bytes())

	// heic 이미지의 exif 데이터 추출
	heifExif, err := goheif.ExtractExif(reader)

	if err != nil {
		log.Println("heic exif error: ", err)
		return nil, nil, err
	}

	return img, heifExif, nil
}
