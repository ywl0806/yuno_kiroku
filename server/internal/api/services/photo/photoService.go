package photo

import (
	"bytes"
	"context"
	"database/sql"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/ywl0806/yuno_kiroku/internal/api/services/photo/models"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	imageHelper "github.com/ywl0806/yuno_kiroku/pkg/imageHelper"
	"github.com/ywl0806/yuno_kiroku/pkg/storage"
)

type PhotoService struct {
	queries         *db.Queries
	standardStorage storage.StorageService
	longTermStorage storage.StorageService
}

func NewPhotoService(
	queries *db.Queries,
	standardStorage storage.StorageService,
	longTermStorage storage.StorageService,
) *PhotoService {
	return &PhotoService{
		queries:         queries,
		standardStorage: standardStorage,
		longTermStorage: longTermStorage,
	}
}

type UploadPhotoReturn struct {
	ThumbnailUrl   string                 `json:"thumbnailUrl"`
	OriginalUrl    string                 `json:"originalUrl"`
	FileName       string                 `json:"fileName"`
	Width          int32                  `json:"width"`
	Height         int32                  `json:"height"`
	Orientation    int32                  `json:"orientation"`
	PhotoCreatedAt time.Time              `json:"photoCreatedAt"`
	FaceDetections []models.FaceDetection `json:"faceDetections"`
}

func (s *PhotoService) CreatePhoto(ctx context.Context, params db.CreatePhotoParams) (*db.Photo, error) {
	photo, err := s.queries.CreatePhoto(ctx, params)
	if err != nil {
		return nil, err
	}
	return &photo, nil
}

/*
*

	Photo 업로드
	1. 파일을 업로드하고 썸네일 이미지를 생성합니다.
	2. 파일 확장자를 확인하고 이미지 파일인 경우 썸네일 이미지를 생성합니다.
	3. 썸네일 이미지를 생성하고 원본 이미지를 저장합니다.
	4. 썸네일 이미지와 원본 이미지의 URL을 반환합니다.
	5. 얼굴 인식 결과를 반환합니다.
*/
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
	err = imgHandler.ResizeImage(1500)
	if err != nil {
		log.Println("resize error: ", err)
		return nil, err
	}
	// resizedFile의 데이터를 새로운 Reader로 만들어서 전달 (버퍼가 이미 읽혔을 수 있으므로)
	faceDetections, err := GetFaceDetection(bytes.NewReader(resizedFile.Bytes()))
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

	originalUrl, err := s.longTermStorage.SaveFile(imgHandler.OriginalFile, folderName, file.Filename)
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
		FaceDetections: *faceDetections,
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

func (s *PhotoService) FindPhotosByPhotoCreatedAt(ctx context.Context, params *db.FindPhotosByPhotoCreatedAtParams) ([]db.Photo, error) {
	photos, err := s.queries.FindPhotosByPhotoCreatedAt(ctx, *params)
	if err != nil {
		return nil, err
	}
	return photos, nil
}

// func (s *PhotoService) FindPhotosGroupByDate(ctx context.Context, params *db.FindPhotosByPhotoCreatedAtParams) ([]db.PhotoGroup, error) {
// 	photos, err := s.FindPhotosByPhotoCreatedAt(ctx, params)

// }

type PhotoGroup struct {
	Year   int        `json:"year" bson:"year"`
	Month  int        `json:"month" bson:"month"`
	Photos []db.Photo `json:"photos" bson:"photos"`
}

func (s *PhotoService) groupPhotosByDate(photos []db.Photo) {

}

type PhotoRange struct {
	Year  string `json:"year"`
	Month string `json:"month"`
}

func (s *PhotoService) GetPhotoRange(ctx context.Context, groupId int32, clanGroupId *int32) ([]PhotoRange, error) {
	// 사진들을 년도와 월로 그룹화
	ranges, err := s.queries.GetPhotoRange(ctx, db.GetPhotoRangeParams{
		GroupID: groupId,
		ClanGroupID: sql.NullInt32{
			Int32: *clanGroupId,
			Valid: clanGroupId != nil,
		},
	})
	if err != nil {
		return nil, err
	}

	photoRanges := make([]PhotoRange, len(ranges))
	for i, r := range ranges {
		photoRanges[i] = PhotoRange{
			Year:  r.Year,
			Month: r.Month,
		}
	}

	return photoRanges, nil
}
