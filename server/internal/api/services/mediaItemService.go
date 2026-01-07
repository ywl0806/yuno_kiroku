package services

import (
	"context"
	"database/sql"
	"log"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/ywl0806/yuno_kiroku/internal/api/consts"
	"github.com/ywl0806/yuno_kiroku/internal/api/enums"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	imageHelper "github.com/ywl0806/yuno_kiroku/pkg/imageHelper"
	"github.com/ywl0806/yuno_kiroku/pkg/storage"
)

type MediaItemService struct {
	queries        *db.Queries
	storageService storage.StorageService
}

func NewMediaItemService(
	queries *db.Queries,
	storageService storage.StorageService,
) *MediaItemService {
	return &MediaItemService{
		queries:        queries,
		storageService: storageService,
	}
}

type CreateMediaItemParams struct {
	ThumbnailStorageKey string
	OriginalStorageKey  string
	ViewStorageKey      string
	ImageHandler        *imageHelper.ImageHelper
	GroupId             int32
	AlbumId             int32
	OriginalFilename    string
}

func (s *MediaItemService) CreateMediaItem(ctx context.Context, params *CreateMediaItemParams) (*db.MediaItem, error) {

	thumbnailFile, _ := params.ImageHandler.GetResizedFile(consts.THUMBNAIL_MAX_LENGTH)
	viewFile, _ := params.ImageHandler.GetResizedFile(consts.VIEW_MAX_LENGTH)

	// 사진 저장 파라미터 생성
	createMediaItemParams := db.CreateMediaItemParams{
		GroupID:  params.GroupId,
		AlbumID:  params.AlbumId,
		TakenAt:  params.ImageHandler.GetTakenAt(),
		FileName: sql.NullString{String: params.OriginalFilename, Valid: true},
	}

	mediaItem, err := s.queries.CreateMediaItem(ctx, createMediaItemParams)
	if err != nil {
		return nil, err
	}

	createOriginalMediaFileParams := db.CreateMediaFileParams{
		MediaItemID: mediaItem.ID,
		Role:        string(enums.MediaItemRoleOriginal),
		StorageKey:  params.OriginalStorageKey,
		// MimeType:    sql.NullString{String: params.ImageHandler.Ext, Valid: true},
		Width:    sql.NullInt32{Int32: int32(params.ImageHandler.OriginalWidth), Valid: true},
		Height:   sql.NullInt32{Int32: int32(params.ImageHandler.OriginalHeight), Valid: true},
		FileSize: sql.NullInt64{Int64: int64(len(params.ImageHandler.OriginalFile)), Valid: true},
	}

	_, err = s.queries.CreateMediaFile(ctx, createOriginalMediaFileParams)
	if err != nil {
		return nil, err
	}

	createThumbnailMediaFileParams := db.CreateMediaFileParams{
		MediaItemID: mediaItem.ID,
		Role:        string(enums.MediaItemRoleThumbnail),
		StorageKey:  params.ThumbnailStorageKey,
		// MimeType:    sql.NullString{String: params.ImageHandler.Ext, Valid: true},
		Width:    sql.NullInt32{Int32: int32(thumbnailFile.Width), Valid: true},
		Height:   sql.NullInt32{Int32: int32(thumbnailFile.Height), Valid: true},
		FileSize: sql.NullInt64{Int64: int64(len(thumbnailFile.File)), Valid: true},
	}

	_, err = s.queries.CreateMediaFile(ctx, createThumbnailMediaFileParams)
	if err != nil {
		return nil, err
	}

	createViewMediaFileParams := db.CreateMediaFileParams{
		MediaItemID: mediaItem.ID,
		Role:        string(enums.MediaItemRoleView),
		StorageKey:  params.ViewStorageKey,
		// MimeType:    sql.NullString{String: params.ImageHandler.Ext, Valid: true},
		Width:    sql.NullInt32{Int32: int32(viewFile.Width), Valid: true},
		Height:   sql.NullInt32{Int32: int32(viewFile.Height), Valid: true},
		FileSize: sql.NullInt64{Int64: int64(len(viewFile.File)), Valid: true},
	}

	_, err = s.queries.CreateMediaFile(ctx, createViewMediaFileParams)
	if err != nil {
		return nil, err
	}

	return &mediaItem, nil
}

func (s *MediaItemService) HandleImage(file *multipart.FileHeader) (*imageHelper.ImageHelper, error) {
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

func (s *MediaItemService) UploadImage(imageHandler *imageHelper.ImageHelper, uploadPath string) (string, string, string, error) {
	takenAt := imageHandler.GetTakenAt()

	folderName := uploadPath + "/" + takenAt.Format("2006-01-02")
	filename := uuid.New().String()
	thumbnailFile, err := imageHandler.GetResizedFile(consts.THUMBNAIL_MAX_LENGTH)
	if err != nil {
		log.Println("Thumbnail File Error: ", err)
		return "", "", "", err
	}
	thumbnailStorageKey, err := s.storageService.SaveFile(thumbnailFile.File, folderName, consts.THUMBNAIL_STORAGE_PREFIX+"/"+filename+thumbnailFile.Ext)
	if err != nil {
		log.Println("Thumbnail Storage Error: ", err)
		return "", "", "", err
	}

	viewFile, err := imageHandler.GetResizedFile(consts.VIEW_MAX_LENGTH)
	if err != nil {
		log.Println("View File Error: ", err)
		return "", "", "", err
	}

	viewStorageKey, err := s.storageService.SaveFile(viewFile.File, folderName, consts.VIEW_STORAGE_PREFIX+"/"+filename+viewFile.Ext)
	if err != nil {
		log.Println("View Storage Error: ", err)
		return "", "", "", err
	}

	originalStorageKey, err := s.storageService.SaveFile(imageHandler.OriginalFile, folderName, consts.ORIGINAL_STORAGE_PREFIX+"/"+filename+"."+imageHandler.Ext)
	if err != nil {
		log.Println("Original Storage Error: ", err)
		return "", "", "", err
	}

	return thumbnailStorageKey, viewStorageKey, originalStorageKey, nil
}

// func (s *PhotoService) UploadLiveMovie(liveMovie *multipart.FileHeader) (string, string, error) {

// 	live, err := liveMovie.Open()

// 	if err != nil {
// 		log.Println("live movie file open error: ", err)
// 		return "", "", err
// 	}
// 	defer live.Close()

// 	liveBytes, err := io.ReadAll(live)
// 	if err != nil {
// 		log.Println("live movie file read error: ", err)
// 		return "", "", err
// 	}
// 	// todo: resize live movie

// 	url, err := s.thumbnailStorage.SaveFile(liveBytes, "live", liveMovie.Filename)
// 	if err != nil {
// 		log.Println("Thumbnail Storage Error : ", err)
// 		return "", "", err
// 	}

// 	return url, url, nil
// }

func (s *MediaItemService) CreateUploadPath(groupId int32, albumId int32) string {
	uploadPath := strconv.Itoa(int(groupId))
	if albumId != 0 {
		uploadPath += "/" + strconv.Itoa(int(albumId))
	}
	return uploadPath
}

func (s *MediaItemService) GetMediaItemsByTakenAt(ctx context.Context, params *db.GetMediaItemsByTakenAtParams) ([]db.GetMediaItemsByTakenAtRow, error) {
	mediaItems, err := s.queries.GetMediaItemsByTakenAt(ctx, *params)
	if err != nil {
		log.Println("get media items by taken at error: ", err)
		return nil, err
	}
	return mediaItems, nil
}

type MediaItemRange struct {
	Year  int `json:"year"`
	Month int `json:"month"`
}

func (s *MediaItemService) GetMediaItemRange(ctx context.Context, clanGroupId int32) ([]MediaItemRange, error) {
	// 사진들을 년도와 월로 그룹화
	ranges, err := s.queries.GetMediaItemRange(ctx, clanGroupId)
	if err != nil {
		return nil, err
	}

	mediaItemRanges := make([]MediaItemRange, len(ranges))

	for i, r := range ranges {
		year, err := strconv.Atoi(r.Year)
		if err != nil {
			return nil, err
		}
		month, err := strconv.Atoi(r.Month)
		if err != nil {
			return nil, err
		}
		mediaItemRanges[i] = MediaItemRange{
			Year:  year,
			Month: month,
		}
	}

	return mediaItemRanges, nil
}
