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
	faceService    *FaceService
}

func NewMediaItemService(
	queries *db.Queries,
	storageService storage.StorageService,
	faceService *FaceService,
) *MediaItemService {
	return &MediaItemService{
		queries:        queries,
		storageService: storageService,
		faceService:    faceService,
	}
}

// 이미지 업로드 결과
type UploadImageResult struct {
	MediaItemID int32
	Status      string
}

// 미디어 아이템 생성 파라미터
type CreateMediaItemParams struct {
	ThumbnailStorageKey string
	OriginalStorageKey  string
	ViewStorageKey      string
	ImageHandler        *imageHelper.ImageHelper
	GroupId             int32
	AlbumId             int32
	OriginalFilename    string
	UploadBatchID       int32
}

// 이미지 처리 백그라운드 처리 파라미터
type ProcessImageParams struct {
	GroupId       int32
	AlbumId       int32
	UploadBatchID int32
	MediaItemID   int32
	ImageHandler  *imageHelper.ImageHelper
	UploadPath    string
	Retry         bool
}

// 미디어 아이템 생성
func (s *MediaItemService) CreateMediaItem(ctx context.Context, params *CreateMediaItemParams) (*db.MediaItem, error) {

	thumbnailFile, _ := params.ImageHandler.GetResizedFile(consts.THUMBNAIL_MAX_LENGTH)
	viewFile, _ := params.ImageHandler.GetResizedFile(consts.VIEW_MAX_LENGTH)

	// 사진 저장 파라미터 생성
	createMediaItemParams := db.CreateMediaItemParams{
		GroupID:       params.GroupId,
		AlbumID:       params.AlbumId,
		UploadBatchID: params.UploadBatchID,
		TakenAt:       params.ImageHandler.GetTakenAt(),
		FileName:      sql.NullString{String: params.OriginalFilename, Valid: true},
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

// 미디어 아이템 생성 (원본 파일만 저장)
func (s *MediaItemService) CreateMediaItemWithOriginalOnly(ctx context.Context, params *CreateMediaItemParams) (*db.MediaItem, error) {
	// 사진 저장 파라미터 생성
	createMediaItemParams := db.CreateMediaItemParams{
		GroupID:       params.GroupId,
		AlbumID:       params.AlbumId,
		UploadBatchID: params.UploadBatchID,
		TakenAt:       params.ImageHandler.GetTakenAt(),
		FileName:      sql.NullString{String: params.OriginalFilename, Valid: true},
	}

	mediaItem, err := s.queries.CreateMediaItem(ctx, createMediaItemParams)
	if err != nil {
		return nil, err
	}

	// 원본 파일만 저장
	createOriginalMediaFileParams := db.CreateMediaFileParams{
		MediaItemID: mediaItem.ID,
		Role:        string(enums.MediaItemRoleOriginal),
		StorageKey:  params.OriginalStorageKey,
		Width:       sql.NullInt32{Int32: int32(params.ImageHandler.OriginalWidth), Valid: true},
		Height:      sql.NullInt32{Int32: int32(params.ImageHandler.OriginalHeight), Valid: true},
		FileSize:    sql.NullInt64{Int64: int64(len(params.ImageHandler.OriginalFile)), Valid: true},
	}

	_, err = s.queries.CreateMediaFile(ctx, createOriginalMediaFileParams)
	if err != nil {
		return nil, err
	}

	return &mediaItem, nil
}

// 이미지 업로드
func (s *MediaItemService) UploadImage(
	ctx context.Context,
	file *multipart.FileHeader,
	groupId int32,
	albumId int32,
	uploadBatchID int32,
	retry bool,
) (*UploadImageResult, error) {
	imgHandler, err := s.HandleImage(file)
	if err != nil {
		return nil, err
	}

	uploadPath := s.CreateUploadPath(groupId, albumId)

	originalStorageKey, err := s.UploadOriginalImage(imgHandler, uploadPath)
	if err != nil {
		return nil, err
	}

	createMediaItemParams := CreateMediaItemParams{
		OriginalStorageKey: originalStorageKey,
		ImageHandler:       imgHandler,
		GroupId:            groupId,
		AlbumId:            albumId,
		OriginalFilename:   file.Filename,
		UploadBatchID:      uploadBatchID,
	}

	mediaItem, err := s.CreateMediaItemWithOriginalOnly(ctx, &createMediaItemParams)
	if err != nil {
		return nil, err
	}

	go func() {
		bgCtx := context.Background()
		s.processImageInBackground(bgCtx, &ProcessImageParams{
			GroupId:       groupId,
			AlbumId:       albumId,
			UploadBatchID: uploadBatchID,
			MediaItemID:   mediaItem.ID,
			ImageHandler:  imgHandler,
			UploadPath:    uploadPath,
			Retry:         retry,
		})
	}()

	return &UploadImageResult{MediaItemID: mediaItem.ID, Status: "uploaded"}, nil
}

// 이미지 처리 백그라운드 처리
func (s *MediaItemService) processImageInBackground(ctx context.Context, p *ProcessImageParams) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("백그라운드 처리 중 패닉 발생: ", r)
			s.UpdateMediaItemUploadStatus(ctx, p.MediaItemID, enums.UploadStatusFailed)
		}
	}()

	viewImageFile, err := p.ImageHandler.GetResizedFile(consts.VIEW_MAX_LENGTH)
	if err != nil {
		log.Println("뷰 이미지 생성 실패: ", err)
		s.UpdateMediaItemUploadStatus(ctx, p.MediaItemID, enums.UploadStatusFailed)
		return
	}

	faceDetections, err := s.faceService.GetFaceDetection(viewImageFile.File)
	if err != nil {
		log.Println("얼굴 인식 실패: ", err)
	}

	if !p.Retry && faceDetections != nil {
		err = s.faceService.CheckImageDuplicateByFaceDetection(ctx, p.GroupId, p.AlbumId, *faceDetections)
		if err != nil {
			log.Println("중복 사진 발견: ", err)
			s.UpdateMediaItemUploadStatus(ctx, p.MediaItemID, enums.UploadStatusDuplicate)
			return
		}
	}

	thumbnailStorageKey, viewStorageKey, err := s.UploadResizedImages(p.ImageHandler, p.UploadPath)
	if err != nil {
		log.Println("리사이즈 이미지 업로드 실패: ", err)
		s.UpdateMediaItemUploadStatus(ctx, p.MediaItemID, enums.UploadStatusFailed)
		return
	}

	err = s.UpdateMediaItemWithResizedImages(ctx, p.MediaItemID, thumbnailStorageKey, viewStorageKey, p.ImageHandler)
	if err != nil {
		log.Println("미디어 파일 업데이트 실패: ", err)
		s.UpdateMediaItemUploadStatus(ctx, p.MediaItemID, enums.UploadStatusFailed)
		return
	}

	if faceDetections != nil {
		_, err = s.faceService.SearchAndSaveFaceDetections(ctx, p.GroupId, p.MediaItemID, *faceDetections)
		if err != nil {
			log.Println("얼굴 인식 결과 저장 실패: ", err)
		}
	}

	s.UpdateMediaItemUploadStatus(ctx, p.MediaItemID, enums.UploadStatusCompleted)
}

// 미디어 아이템 업데이트 (썸네일 파일과 뷰 파일 추가)
func (s *MediaItemService) UpdateMediaItemWithResizedImages(
	ctx context.Context,
	mediaItemID int32,
	thumbnailStorageKey string,
	viewStorageKey string,
	imageHandler *imageHelper.ImageHelper,
) error {
	thumbnailFile, _ := imageHandler.GetResizedFile(consts.THUMBNAIL_MAX_LENGTH)
	viewFile, _ := imageHandler.GetResizedFile(consts.VIEW_MAX_LENGTH)

	// 썸네일 파일 저장
	createThumbnailMediaFileParams := db.CreateMediaFileParams{
		MediaItemID: mediaItemID,
		Role:        string(enums.MediaItemRoleThumbnail),
		StorageKey:  thumbnailStorageKey,
		Width:       sql.NullInt32{Int32: int32(thumbnailFile.Width), Valid: true},
		Height:      sql.NullInt32{Int32: int32(thumbnailFile.Height), Valid: true},
		FileSize:    sql.NullInt64{Int64: int64(len(thumbnailFile.File)), Valid: true},
	}

	_, err := s.queries.CreateMediaFile(ctx, createThumbnailMediaFileParams)
	if err != nil {
		return err
	}

	// 뷰 파일 저장
	createViewMediaFileParams := db.CreateMediaFileParams{
		MediaItemID: mediaItemID,
		Role:        string(enums.MediaItemRoleView),
		StorageKey:  viewStorageKey,
		Width:       sql.NullInt32{Int32: int32(viewFile.Width), Valid: true},
		Height:      sql.NullInt32{Int32: int32(viewFile.Height), Valid: true},
		FileSize:    sql.NullInt64{Int64: int64(len(viewFile.File)), Valid: true},
	}

	_, err = s.queries.CreateMediaFile(ctx, createViewMediaFileParams)
	if err != nil {
		return err
	}

	return nil
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

// UploadImageAll uploads thumbnail, view, and original images in one go (legacy helper).
func (s *MediaItemService) UploadImageAll(imageHandler *imageHelper.ImageHelper, uploadPath string) (string, string, string, error) {
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

// 원본 이미지 업로드
func (s *MediaItemService) UploadOriginalImage(imageHandler *imageHelper.ImageHelper, uploadPath string) (string, error) {
	takenAt := imageHandler.GetTakenAt()
	folderName := uploadPath + "/" + takenAt.Format("2006-01-02")
	filename := uuid.New().String()

	originalStorageKey, err := s.storageService.SaveFile(imageHandler.OriginalFile, folderName, consts.ORIGINAL_STORAGE_PREFIX+"/"+filename+"."+imageHandler.Ext)
	if err != nil {
		log.Println("Original Storage Error: ", err)
		return "", err
	}

	return originalStorageKey, nil
}

// 리사이즈 이미지 업로드
func (s *MediaItemService) UploadResizedImages(imageHandler *imageHelper.ImageHelper, uploadPath string) (string, string, error) {
	takenAt := imageHandler.GetTakenAt()
	folderName := uploadPath + "/" + takenAt.Format("2006-01-02")
	filename := uuid.New().String()

	thumbnailFile, err := imageHandler.GetResizedFile(consts.THUMBNAIL_MAX_LENGTH)
	if err != nil {
		log.Println("Thumbnail File Error: ", err)
		return "", "", err
	}
	thumbnailStorageKey, err := s.storageService.SaveFile(thumbnailFile.File, folderName, consts.THUMBNAIL_STORAGE_PREFIX+"/"+filename+thumbnailFile.Ext)
	if err != nil {
		log.Println("Thumbnail Storage Error: ", err)
		return "", "", err
	}

	viewFile, err := imageHandler.GetResizedFile(consts.VIEW_MAX_LENGTH)
	if err != nil {
		log.Println("View File Error: ", err)
		return "", "", err
	}

	viewStorageKey, err := s.storageService.SaveFile(viewFile.File, folderName, consts.VIEW_STORAGE_PREFIX+"/"+filename+viewFile.Ext)
	if err != nil {
		log.Println("View Storage Error: ", err)
		return "", "", err
	}

	return thumbnailStorageKey, viewStorageKey, nil
}

// 업로드 경로 생성
func (s *MediaItemService) CreateUploadPath(groupId int32, albumId int32) string {
	uploadPath := strconv.Itoa(int(groupId))
	if albumId != 0 {
		uploadPath += "/" + strconv.Itoa(int(albumId))
	}
	return uploadPath
}

// 미디어 아이템 조회
func (s *MediaItemService) GetMediaItemsByTakenAt(ctx context.Context, params *db.GetMediaItemsByTakenAtParams) ([]db.GetMediaItemsByTakenAtRow, error) {
	mediaItems, err := s.queries.GetMediaItemsByTakenAt(ctx, *params)
	if err != nil {
		log.Println("get media items by taken at error: ", err)
		return nil, err
	}
	return mediaItems, nil
}

// 미디어 아이템 범위
type MediaItemRange struct {
	Year  int `json:"year"`
	Month int `json:"month"`
}

// 미디어 아이템 범위 조회
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

// 업로드 배치 생성
func (s *MediaItemService) CreateUploadBatch(ctx context.Context, groupId int32, albumId int32) (*db.UploadBatch, error) {
	uploadBatch, err := s.queries.CreateUploadBatch(ctx, albumId)
	if err != nil {
		return nil, err
	}
	return &uploadBatch, nil
}

// 미디어 아이템 업로드 상태 업데이트
func (s *MediaItemService) UpdateMediaItemUploadStatus(ctx context.Context, mediaItemID int32, uploadStatus enums.UploadStatus) error {
	_, err := s.queries.UpdateMediaItemUploadStatus(ctx, db.UpdateMediaItemUploadStatusParams{
		ID:           mediaItemID,
		UploadStatus: string(uploadStatus),
	})
	return err
}

// 업로드 배치 상태 조회
func (s *MediaItemService) GetUploadStatuses(ctx context.Context, uploadBatchID int32) ([]db.GetUploadStatusesRow, error) {
	uploadStatuses, err := s.queries.GetUploadStatuses(ctx, uploadBatchID)
	if err != nil {
		return nil, err
	}
	return uploadStatuses, nil
}
