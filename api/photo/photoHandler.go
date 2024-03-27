package photo

import (
	"bytes"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"github.com/ywl0806/yuno_kiroku/api/lib/imageHandler"
	"github.com/ywl0806/yuno_kiroku/api/utils"
)

type UploadPhotoReturn struct {
	ThumbnailUrl   string    `json:"thumbnailUrl"`
	OriginalUrl    string    `json:"originalUrl"`
	FileName       string    `json:"fileName"`
	PhotoCreatedAt time.Time `json:"photoCreatedAt"`
}

func (con *PhotoController) uploadPhoto(file *multipart.FileHeader) (*UploadPhotoReturn, error) {

	ext := strings.Split(file.Filename, ".")[1]
	originalFile, err := file.Open()
	if err != nil {
		log.Println("file open error: ", err)
		return nil, err
	}
	defer originalFile.Close()

	buf1, buf2, _ := utils.CopyReader(originalFile)

	// file resize
	// convert to jpeg
	// get exif
	resizedFile := bytes.NewBuffer(nil)
	exifData, err := imageHandler.ResizeImage(buf1, resizedFile, ext)

	if err != nil {
		log.Println("resize error: ", err)
		return nil, err
	}
	photoCreatedAt, _ := exifData.DateTime()

	if photoCreatedAt.IsZero() {
		photoCreatedAt = time.Now()
	}

	originalFilename := strings.Split(file.Filename, ".")[0]

	var folderName string

	now := time.Now()
	folderName = now.Format("2006-01-02")

	thumbnailUrl, err := con.standardStorage.SaveFile(resizedFile, folderName, originalFilename+".jpeg")
	if err != nil {
		log.Println("Standard Storage Error: ", err)
		return nil, err
	}

	originalUrl, err := con.longTermStorage.SaveFile(buf2, "", file.Filename)
	if err != nil {
		log.Println("Longterm Storage Error: ", err)
		return nil, err
	}

	result := UploadPhotoReturn{
		ThumbnailUrl:   thumbnailUrl,
		OriginalUrl:    originalUrl,
		FileName:       file.Filename,
		PhotoCreatedAt: photoCreatedAt,
	}

	return &result, nil
}

type UploadLiveMovieReturn struct {
	LiveUrl         string `json:"liveUrl"`
	OriginalLiveUrl string `json:"originalLiveUrl"`
}

func (con *PhotoController) uploadLiveMovie(liveMovie *multipart.FileHeader) (*UploadLiveMovieReturn, error) {

	live, err := liveMovie.Open()

	if err != nil {
		log.Println("live movie file open error: ", err)
		return nil, err
	}

	// todo: resize live movie

	url, err := con.standardStorage.SaveFile(live, "live", liveMovie.Filename)
	if err != nil {
		log.Println("Standard Storage Error : ", err)
		return nil, err
	}

	result := UploadLiveMovieReturn{
		LiveUrl:         url,
		OriginalLiveUrl: url,
	}

	return &result, nil
}
