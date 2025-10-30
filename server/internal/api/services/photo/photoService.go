package photo

import (
	"bytes"
	"context"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/ywl0806/yuno_kiroku/internal/db"
	imageHelper "github.com/ywl0806/yuno_kiroku/pkg/imageHelper"
	"github.com/ywl0806/yuno_kiroku/pkg/storage"
)

type PhotoService struct {
	queries         *db.Queries
	standardStorage storage.StorageService
	longTermStorage storage.StorageService
}

func NewPhotoService(standardStorage storage.StorageService, longTermStorage storage.StorageService) *PhotoService {
	return &PhotoService{
		standardStorage: standardStorage,
		longTermStorage: longTermStorage,
	}
}

type UploadPhotoReturn struct {
	ThumbnailUrl   string    `json:"thumbnailUrl"`
	OriginalUrl    string    `json:"originalUrl"`
	FileName       string    `json:"fileName"`
	Width          int32     `json:"width"`
	Height         int32     `json:"height"`
	Orientation    int32     `json:"orientation"`
	PhotoCreatedAt time.Time `json:"photoCreatedAt"`
}

func (s *PhotoService) CreatePhoto(ctx context.Context, params db.CreatePhotoParams) (*db.Photo, error) {
	photo, err := s.queries.CreatePhoto(ctx, params)
	if err != nil {
		return nil, err
	}
	return &photo, nil
}

func (s *PhotoService) UploadPhoto(file *multipart.FileHeader, uploadPath string) (*UploadPhotoReturn, error) {

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
	imgHandler := imageHelper.NewImageHandler(originalFile, resizedFile, ext)
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
	folderName = uploadPath + "/" + now.Format("2006-01-02")

	thumbnailUrl, err := s.standardStorage.SaveFile(resizedFile, folderName, originalFilename+".jpeg")
	if err != nil {
		log.Println("Standard Storage Error: ", err)
		return nil, err
	}

	originalUrl, err := s.longTermStorage.SaveFile(imgHandler.OriginalFile, "", file.Filename)
	if err != nil {
		log.Println("Longterm Storage Error: ", err)
		return nil, err
	}
	orientationRaw, _ := imgHandler.Exif.Get("Orientation")
	orientation := 1
	if orientationRaw != nil {
		orientation, err = orientationRaw.Int(0)
	}
	if err != nil {
		orientation = 1
		err = nil
	}

	result := UploadPhotoReturn{
		ThumbnailUrl:   thumbnailUrl,
		OriginalUrl:    originalUrl,
		FileName:       file.Filename,
		Width:          int32(imgHandler.OriginalImage.Bounds().Dx()),
		Height:         int32(imgHandler.OriginalImage.Bounds().Dy()),
		Orientation:    int32(orientation),
		PhotoCreatedAt: photoCreatedAt,
	}

	return &result, nil
}

type UploadLiveMovieReturn struct {
	LiveUrl         string `json:"liveUrl"`
	OriginalLiveUrl string `json:"originalLiveUrl"`
}

func (s *PhotoService) UploadLiveMovie(liveMovie *multipart.FileHeader) (*UploadLiveMovieReturn, error) {

	live, err := liveMovie.Open()

	if err != nil {
		log.Println("live movie file open error: ", err)
		return nil, err
	}

	// todo: resize live movie

	url, err := s.standardStorage.SaveFile(live, "live", liveMovie.Filename)
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

func (s *PhotoService) CreateUploadPath(groupId int32, clanGroupId *int32) string {
	uploadPath := strconv.Itoa(int(groupId))
	if clanGroupId != nil {
		uploadPath += "/" + strconv.Itoa(int(*clanGroupId))
	}
	return uploadPath
}
