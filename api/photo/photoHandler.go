package photo

import (
	"bytes"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"github.com/ywl0806/yuno_kiroku/api/lib/imageHandler"
)

type UploadPhotoReturn struct {
	ThumbnailUrl   string    `json:"thumbnailUrl"`
	OriginalUrl    string    `json:"originalUrl"`
	FileName       string    `json:"fileName"`
	Width          int       `json:"width"`
	Height         int       `json:"height"`
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

	// file resize
	// convert to jpeg
	// get exif
	resizedFile := new(bytes.Buffer)
	imgHandler := imageHandler.NewImageHandler(originalFile, resizedFile, ext)
	err = imgHandler.ResizeImage(1500, 1500)

	if err != nil {
		log.Println("resize error: ", err)
		return nil, err
	}

	photoCreatedAt, _ := imgHandler.Exif.DateTime()

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

	originalUrl, err := con.longTermStorage.SaveFile(imgHandler.OriginalFile, "", file.Filename)
	if err != nil {
		log.Println("Longterm Storage Error: ", err)
		return nil, err
	}

	result := UploadPhotoReturn{
		ThumbnailUrl:   thumbnailUrl,
		OriginalUrl:    originalUrl,
		FileName:       file.Filename,
		Width:          imgHandler.OriginalImage.Bounds().Dx(),
		Height:         imgHandler.OriginalImage.Bounds().Dy(),
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
