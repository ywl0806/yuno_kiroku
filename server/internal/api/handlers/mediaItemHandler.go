package handlers

import (
	"context"
	"log"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/ywl0806/yuno_kiroku/internal/api/consts"
	"github.com/ywl0806/yuno_kiroku/internal/api/enums"
	customErrors "github.com/ywl0806/yuno_kiroku/internal/api/errors"
	"github.com/ywl0806/yuno_kiroku/internal/api/handlers/models"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
	"github.com/ywl0806/yuno_kiroku/internal/api/services"
	"github.com/ywl0806/yuno_kiroku/internal/api/utils"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	imageHelper "github.com/ywl0806/yuno_kiroku/pkg/imageHelper"
)

type MediaItemHandler struct {
	mediaItemService *services.MediaItemService
	faceService      *services.FaceService
}

func NewMediaItemHandler(
	mediaItemService *services.MediaItemService,
	faceService *services.FaceService,
) *MediaItemHandler {
	return &MediaItemHandler{
		mediaItemService: mediaItemService,
		faceService:      faceService,
	}
}

// @Tags MediaItem
// @Description 업로드 배치 생성
// @Param album_id query string true "Album ID"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /media-item/upload-batch [post]
func (con *MediaItemHandler) CreateUploadBatch(c echo.Context) error {
	// album 쿼리 파라미터 가져오기
	albumId, err := utils.ConvertToInt32(c.QueryParam("album_id"))
	if err != nil {
		return customErrors.NewRequiredError(consts.Album)
	}

	// 그룹 ID 가져오기
	authUser := middlewares.GetAuthUser(c)
	groupId := authUser.GroupId

	// 컨텍스트 가져오기
	ctx := c.Request().Context()

	// 업로드 배치 생성
	uploadBatch, err := con.mediaItemService.CreateUploadBatch(ctx, groupId, albumId)
	if err != nil {
		return err
	}

	return c.JSON(200, map[string]interface{}{
		"id":         uploadBatch.ID,
		"album_id":   uploadBatch.AlbumID,
		"created_at": uploadBatch.CreatedAt,
	})
}

// @Tags MediaItem
// @Description 이미지 업로드
// @Accept  multipart/form-data
// @Param file formData file true "file"
// @Param album_id query string true "Album ID"
// @Param upload_batch_id query string true "Upload Batch ID"
// @Param retry query string false "Retry" example(1)
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /media-item/image/upload [post]
func (con *MediaItemHandler) UploadImage(c echo.Context) error {
	// 파일 가져오기
	file, err := c.FormFile("file")
	if err != nil {
		return customErrors.NewRequiredError(consts.Photo)
	}

	// album 쿼리 파라미터 가져오기
	albumId, err := utils.ConvertToInt32(c.QueryParam("album_id"))
	if err != nil {
		return customErrors.NewRequiredError(consts.Album)
	}

	// upload_batch_id 쿼리 파라미터 가져오기 (필수)
	uploadBatchID, err := utils.ConvertToInt32(c.QueryParam("upload_batch_id"))
	if err != nil {
		return customErrors.NewRequiredError("upload_batch_id")
	}

	retry, _ := utils.ConvertToBool(c.QueryParam("retry"))

	// 그룹 ID 가져오기
	authUser := middlewares.GetAuthUser(c)
	groupId := authUser.GroupId

	// 컨텍스트 가져오기
	ctx := c.Request().Context()

	// 이미지 핸들러 생성
	imgHandler, err := con.mediaItemService.HandleImage(file)
	if err != nil {
		log.Println("handle image error: ", err)
		return err
	}

	// 업로드 경로 생성
	uploadPath := con.mediaItemService.CreateUploadPath(groupId, albumId)

	// 원본 파일만 먼저 저장 (즉시 반환을 위해)
	originalStorageKey, err := con.mediaItemService.UploadOriginalImage(imgHandler, uploadPath)
	if err != nil {
		log.Println("upload original image error: ", err)
		return err
	}

	// 미디어 아이템 생성 (원본만 저장)
	createMediaItemParams := services.CreateMediaItemParams{
		OriginalStorageKey: originalStorageKey,

		ImageHandler:     imgHandler,
		GroupId:          groupId,
		AlbumId:          albumId,
		OriginalFilename: file.Filename,
		UploadBatchID:    uploadBatchID,
	}

	mediaItem, err := con.mediaItemService.CreateMediaItemWithOriginalOnly(ctx, &createMediaItemParams)
	if err != nil {
		log.Println("create media item with original only error: ", err)
		return err
	}

	// 백그라운드에서 리사이즈, 얼굴인식 등 처리
	go func() {
		bgCtx := context.Background()
		con.processImageInBackground(bgCtx, groupId, albumId, uploadBatchID, mediaItem.ID, imgHandler, uploadPath, retry)
	}()

	return c.JSON(200, map[string]interface{}{
		"media_item_id": mediaItem.ID,
		"status":        "uploaded",
	})
}

// processImageInBackground processes image resizing and face detection in background
func (con *MediaItemHandler) processImageInBackground(
	ctx context.Context,
	groupId int32,
	albumId int32,
	uploadBatchID int32,
	mediaItemID int32,
	imgHandler *imageHelper.ImageHelper,
	uploadPath string,
	retry bool,
) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("백그라운드 처리 중 패닉 발생: ", r)
			con.mediaItemService.UpdateMediaItemUploadStatus(ctx, mediaItemID, "04")
		}
	}()

	// 뷰 이미지 파일 생성
	viewImageFile, err := imgHandler.GetResizedFile(consts.VIEW_MAX_LENGTH)
	if err != nil {
		log.Println("뷰 이미지 생성 실패: ", err)
		con.mediaItemService.UpdateMediaItemUploadStatus(ctx, mediaItemID, enums.UploadStatusFailed)
		return
	}

	// 얼굴 인식 결과 가져오기
	faceDetections, err := con.faceService.GetFaceDetection(viewImageFile.File)
	if err != nil {
		log.Println("얼굴 인식 실패: ", err)
		// 얼굴 인식 실패해도 계속 진행
	}

	// 사진 중복체크(재업로드가 아닌 경우만 체크)
	if !retry && faceDetections != nil {
		err = con.faceService.CheckImageDuplicateByFaceDetection(ctx, groupId, albumId, *faceDetections)
		if err != nil {
			log.Println("중복 사진 발견: ", err)
			con.mediaItemService.UpdateMediaItemUploadStatus(ctx, mediaItemID, enums.UploadStatusDuplicate)
			return
		}
	}

	// 썸네일과 뷰 이미지 업로드
	thumbnailStorageKey, viewStorageKey, err := con.mediaItemService.UploadResizedImages(imgHandler, uploadPath)
	if err != nil {
		log.Println("리사이즈 이미지 업로드 실패: ", err)
		con.mediaItemService.UpdateMediaItemUploadStatus(ctx, mediaItemID, enums.UploadStatusFailed)
		return
	}

	// 미디어 파일 업데이트 (썸네일, 뷰 추가)
	err = con.mediaItemService.UpdateMediaItemWithResizedImages(ctx, mediaItemID, thumbnailStorageKey, viewStorageKey, imgHandler)
	if err != nil {
		log.Println("미디어 파일 업데이트 실패: ", err)
		con.mediaItemService.UpdateMediaItemUploadStatus(ctx, mediaItemID, enums.UploadStatusFailed)
		return
	}

	// 얼굴 인식 결과 저장
	if faceDetections != nil {
		_, err = con.faceService.SearchAndSaveFaceDetections(ctx, groupId, mediaItemID, *faceDetections)
		if err != nil {
			log.Println("얼굴 인식 결과 저장 실패: ", err)
			// 얼굴 인식 실패해도 사진은 저장되었으므로 계속 진행
		}
	}

	// 배치 상태를 completed로 업데이트
	con.mediaItemService.UpdateMediaItemUploadStatus(ctx, mediaItemID, enums.UploadStatusCompleted)
}

// @Tags MediaItem
// @Description 사진이 있는 년도와 월 목록 조회
// @Router /media-item/range [get]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Success 200 {object} []services.MediaItemRange
func (con *MediaItemHandler) GetMediaItemRange(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)
	ranges, err := con.mediaItemService.GetMediaItemRange(c.Request().Context(), authUser.ClanGroupId)
	if err != nil {
		return err
	}
	return c.JSON(200, ranges)
}

type GetMediaItemsRequest struct {
	From *time.Time `query:"from" validate:"required"`
	To   *time.Time `query:"to" validate:"required"`
}

func (GetMediaItemsRequest) bind(c echo.Context, p *db.GetMediaItemsByTakenAtParams) error {
	reqParams := new(GetMediaItemsRequest)
	if err := c.Bind(reqParams); err != nil {
		return err
	}
	if err := c.Validate(reqParams); err != nil {
		return err
	}

	authUser := middlewares.GetAuthUser(c)

	p.TakenAtFrom = *reqParams.From
	p.TakenAtTo = *reqParams.To
	p.ClanGroupID = authUser.ClanGroupId
	return nil
}

// @Tags MediaItem
// @Description 미디어 아이템 목록 조회
// @Router /media-item [get]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Param from query string true "From" example(2025-01-01)
// @Param to query string true "To" example(2025-01-01)
// @Success 200
func (con *MediaItemHandler) GetMediaItems(c echo.Context) error {
	reqParams := new(GetMediaItemsRequest)

	params := new(db.GetMediaItemsByTakenAtParams)
	if err := reqParams.bind(c, params); err != nil {
		log.Println("bind get media items request error: ", err)
		return err
	}
	mediaItems, err := con.mediaItemService.GetMediaItemsByTakenAt(c.Request().Context(), params)
	if err != nil {
		log.Println("get media items by taken at error: ", err)
		return err
	}

	return c.JSON(200, models.NewMediaItemsResponse(mediaItems))

}

// @Tags MediaItem
// @Description 업로드 배치 상태 조회
// @Param upload_batch_id query string true "Upload Batch ID"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /media-item/upload-batch/status [get]
func (con *MediaItemHandler) GetUploadBatchStatus(c echo.Context) error {
	// upload_batch_id 쿼리 파라미터 가져오기
	uploadBatchID, err := utils.ConvertToInt32(c.QueryParam("upload_batch_id"))
	if err != nil {
		return customErrors.NewRequiredError("upload_batch_id")
	}

	// 컨텍스트 가져오기
	ctx := c.Request().Context()

	// 업로드 배치 조회
	uploadStatuses, err := con.mediaItemService.GetUploadStatuses(ctx, uploadBatchID)
	if err != nil {
		return err
	}

	uploadStatusesResponse := models.NewUploadBatchStatusResponse(uploadStatuses)

	return c.JSON(200, uploadStatusesResponse)
}

// // @Tags Photo
// // @Description upload live photo
// // @Accept  multipart/form-data
// // @Param photo formData file true "photo"
// // @Param live formData file true "live"
// // @Param album_id query string true "Album ID"
// // @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// // @Router /photo/upload-live [post]
// func (con *PhotoHandler) UploadLivePhoto(c echo.Context) error {
// 	photoFile, err := c.FormFile("photo")
// 	if err != nil {
// 		log.Println("no photo formfile error: ", err)
// 		return customErrors.NewRequiredError(consts.Photo)
// 	}

// 	liveMovie, err := c.FormFile("live")
// 	if err != nil {
// 		log.Println("no live movie formfile error: ", err)
// 		return customErrors.NewRequiredError(consts.LiveMovie)
// 	}
// 	groupId := middlewares.GetAuthUser(c).GroupId

// 	albumId, err := utils.ConvertToInt32(c.QueryParam("album_id"))
// 	if err != nil {
// 		return customErrors.NewRequiredError(consts.Album)
// 	}

// 	// 컨텍스트 가져오기
// 	ctx := c.Request().Context()

// 	// 이미지 핸들러 생성(이미지 리사이즈)
// 	imgHandler, err := con.mediaItemService.HandleImage(photoFile)
// 	if err != nil {
// 		return err
// 	}
// 	// 얼굴 인식 결과 가져오기(이미지 리사이즈 된 이미지 이용)
// 	faceDetections, err := con.faceService.GetFaceDetection(imgHandler.ResizedFile)
// 	if err != nil {
// 		return err
// 	}

// 	// 업로드 경로 생성
// 	uploadPath := con.mediaItemService.CreateUploadPath(groupId, albumId)

// 	// 썸네일 이미지와 원본 이미지 업로드
// 	thumbnailUrl, originalUrl, err := con.mediaItemService.UploadPhoto(imgHandler, uploadPath)
// 	if err != nil {
// 		return err
// 	}

// 	// 라이브 무비 업로드
// 	liveUrl, originalLiveUrl, err := con.mediaItemService.UploadLiveMovie(liveMovie)
// 	if err != nil {
// 		log.Println("upload live movie error: ", err)
// 		return err
// 	}

// 	// 사진 저장 파라미터 생성
// 	createPhotoParams := services.CreatePhotoParams{
// 		ThumbnailUrl:     thumbnailUrl,
// 		OriginalUrl:      originalUrl,
// 		LiveUrl:          liveUrl,
// 		OriginalLiveUrl:  originalLiveUrl,
// 		ImageHandler:     imgHandler,
// 		GroupId:          groupId,
// 		AlbumId:          albumId,
// 		OriginalFilename: photoFile.Filename,
// 	}

// 	// 사진 저장
// 	photo, err := con.mediaItemService.CreatePhoto(ctx, &createPhotoParams)
// 	if err != nil {
// 		return err
// 	}

// 	// 얼굴 인식 결과 저장
// 	_, err = con.faceService.SearchAndSaveFaceDetections(ctx, groupId, photo.ID, *faceDetections)
// 	if err != nil {
// 		log.Println("얼굴 인식 결과 저장 실패: ", err)
// 		// 얼굴 인식 실패해도 사진은 저장되었으므로 계속 진행
// 	}

// 	return c.JSON(200, models.NewPhotoResponse(photo))

// }
