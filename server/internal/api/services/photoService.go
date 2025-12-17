package services

import (
	"context"
	"database/sql"
	"io"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ywl0806/yuno_kiroku/internal/api/consts"
	apiErrors "github.com/ywl0806/yuno_kiroku/internal/api/errors"
	"github.com/ywl0806/yuno_kiroku/internal/api/handlers/models"
	"github.com/ywl0806/yuno_kiroku/internal/api/utils"
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

func (s *PhotoService) GetPhotos(ctx context.Context, params *db.FindPhotosByPhotoCreatedAtParams) ([]db.Photo, error) {

	photos, err := s.queries.FindPhotosByPhotoCreatedAt(ctx, *params)
	if err != nil {
		return nil, err
	}
	return photos, nil
}

type CreatePhotoParams struct {
	ThumbnailUrl     string
	OriginalUrl      string
	LiveUrl          string
	OriginalLiveUrl  string
	ImageHandler     *imageHelper.ImageHelper
	GroupId          int32
	AlbumId          int32
	OriginalFilename string
}

func (s *PhotoService) CreatePhoto(ctx context.Context, params *CreatePhotoParams) (*db.Photo, error) {

	thumbnailWidth, thumbnailHeight := params.ImageHandler.GetResizedImageSize()
	originalWidth, originalHeight := params.ImageHandler.GetOriginalImageSize()

	// 사진 저장 파라미터 생성
	createPhotoParams := db.CreatePhotoParams{
		GroupID:         params.GroupId,
		AlbumID:         params.AlbumId,
		ThumbnailUrl:    params.ThumbnailUrl,
		FileName:        params.OriginalFilename,
		OriginalWidth:   sql.NullInt32{Int32: int32(originalWidth), Valid: true},
		OriginalHeight:  sql.NullInt32{Int32: int32(originalHeight), Valid: true},
		ThumbnailWidth:  int32(thumbnailWidth),
		ThumbnailHeight: int32(thumbnailHeight),
		PhotoCreatedAt:  params.ImageHandler.GetPhotoCreatedAt(),
	}

	if params.OriginalUrl != "" {
		createPhotoParams.OriginalUrl = sql.NullString{String: params.OriginalUrl, Valid: true}
	} else {
		createPhotoParams.OriginalUrl = sql.NullString{Valid: false}
	}

	if params.LiveUrl != "" {
		createPhotoParams.LiveUrl = sql.NullString{String: params.LiveUrl, Valid: true}
	} else {
		createPhotoParams.LiveUrl = sql.NullString{Valid: false}
	}

	if params.OriginalLiveUrl != "" {
		createPhotoParams.OriginalLiveUrl = sql.NullString{String: params.OriginalLiveUrl, Valid: true}
	} else {
		createPhotoParams.OriginalLiveUrl = sql.NullString{Valid: false}
	}

	photo, err := s.queries.CreatePhoto(ctx, createPhotoParams)
	if err != nil {
		return nil, err
	}

	return &photo, nil
}

type UploadPhotoReturn struct {
	ThumbnailUrl    string                 `json:"thumbnailUrl"`
	OriginalUrl     string                 `json:"originalUrl"`
	FileName        string                 `json:"fileName"`
	OriginalWidth   int32                  `json:"originalWidth"`
	OriginalHeight  int32                  `json:"originalHeight"`
	ThumbnailWidth  int32                  `json:"thumbnailWidth"`
	ThumbnailHeight int32                  `json:"thumbnailHeight"`
	PhotoCreatedAt  time.Time              `json:"photoCreatedAt"`
	FaceDetections  []models.FaceDetection `json:"faceDetections"`
}

func (s *PhotoService) HandleImage(file *multipart.FileHeader) (*imageHelper.ImageHelper, error) {
	ext := strings.Split(file.Filename, ".")[1]
	originalFile, err := file.Open()
	if err != nil {
		log.Println("file open error: ", err)
		return nil, err
	}
	defer originalFile.Close()

	imgHandler, err := imageHelper.NewImageHandler(originalFile, ext)
	if err != nil {
		log.Println("image handler error: ", err)
		return nil, err
	}

	return imgHandler, nil
}

func (s *PhotoService) UploadPhoto(imageHandler *imageHelper.ImageHelper, uploadPath string) (string, string, error) {
	photoCreatedAt := imageHandler.GetPhotoCreatedAt()

	folderName := uploadPath + "/" + photoCreatedAt.Format("2006-01-02")
	filename := uuid.New().String()

	thumbnailUrl, err := s.standardStorage.SaveFile(imageHandler.ResizedFile, folderName, filename+".jpeg")
	if err != nil {
		log.Println("Standard Storage Error: ", err)
		return "", "", err
	}

	originalUrl, err := s.longTermStorage.SaveFile(imageHandler.OriginalFile, folderName, filename+"."+imageHandler.Ext)
	if err != nil {
		log.Println("Longterm Storage Error: ", err)
		return "", "", err
	}

	return thumbnailUrl, originalUrl, nil
}

/*
*

	얼굴 인식 결과로 사진이 중복되는지 확인합니다.
	1. 얼굴 인식 결과를 벡터로 변환합니다.
	2. 데이터베이스에서 얼굴 인식 결과와 일치하는 사진을 조회합니다.
	3. 조회된 사진이 있으면 중복된 사진이 있다는 에러를 반환합니다.
*/
func (s *PhotoService) CheckPhotoDuplicateByFaceDetection(ctx context.Context, groupId int32, albumId int32, faceDetections []models.FaceDetection) error {

	embeddings := make([]interface{}, len(faceDetections))
	for i, faceDetection := range faceDetections {
		embeddings[i] = utils.Float64SliceToVectorString(faceDetection.Embedding)
	}
	photo, err := s.queries.GetPhotoByFaceDetection(ctx, db.GetPhotoByFaceDetectionParams{
		GroupID:    groupId,
		Embeddings: embeddings,
	})
	if err == sql.ErrNoRows {
		return nil
	}
	if photo.ID != 0 {
		return apiErrors.NewDuplicateError(consts.Photo)
	}

	return nil
}

func (s *PhotoService) UploadLiveMovie(liveMovie *multipart.FileHeader) (string, string, error) {

	live, err := liveMovie.Open()

	if err != nil {
		log.Println("live movie file open error: ", err)
		return "", "", err
	}
	defer live.Close()

	liveBytes, err := io.ReadAll(live)
	if err != nil {
		log.Println("live movie file read error: ", err)
		return "", "", err
	}
	// todo: resize live movie

	url, err := s.standardStorage.SaveFile(liveBytes, "live", liveMovie.Filename)
	if err != nil {
		log.Println("Standard Storage Error : ", err)
		return "", "", err
	}

	return url, url, nil
}

func (s *PhotoService) CreateUploadPath(groupId int32, albumId int32) string {
	uploadPath := strconv.Itoa(int(groupId))
	if albumId != 0 {
		uploadPath += "/" + strconv.Itoa(int(albumId))
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

type PhotoRange struct {
	Year  int `json:"year"`
	Month int `json:"month"`
}

func (s *PhotoService) GetPhotoRange(ctx context.Context, clanGroupId int32) ([]PhotoRange, error) {
	// 사진들을 년도와 월로 그룹화
	ranges, err := s.queries.GetPhotoRange(ctx, clanGroupId)
	if err != nil {
		return nil, err
	}

	photoRanges := make([]PhotoRange, len(ranges))

	for i, r := range ranges {
		year, err := strconv.Atoi(r.Year)
		if err != nil {
			return nil, err
		}
		month, err := strconv.Atoi(r.Month)
		if err != nil {
			return nil, err
		}
		photoRanges[i] = PhotoRange{
			Year:  year,
			Month: month,
		}
	}

	return photoRanges, nil
}

/*
*

 1. Identity ID를 기반으로 랜덤 사진을 가져옵니다.
 2. 사진 정보와 얼굴 위치정보를 반환합니다.
*/
func (s *PhotoService) GetIdentityRandomPhoto(ctx context.Context, clanGroupId int32, identityId int32) (*db.GetIdentityRandomPhotoRow, error) {
	photo, err := s.queries.GetIdentityRandomPhoto(ctx, db.GetIdentityRandomPhotoParams{
		ClanGroupID: clanGroupId,
		IdentityID:  identityId,
	})
	if err == sql.ErrNoRows {
		return nil, apiErrors.NewNotFoundError("photo")
	}
	if err != nil {
		return nil, err
	}
	return &photo, nil
}
