package handlers

import (
	"log"

	"github.com/labstack/echo/v4"

	"github.com/ywl0806/yuno_kiroku/internal/handlers/models"
	"github.com/ywl0806/yuno_kiroku/internal/middlewares"
	"github.com/ywl0806/yuno_kiroku/internal/services"
	"github.com/ywl0806/yuno_kiroku/pkg/utils"

	"github.com/ywl0806/yuno_kiroku/internal/apperr"
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
		return apperr.NewValidationError("message.validation.required", map[string]string{"field": "message-item.album_id"})
	}

	authUser := middlewares.GetAuthUser(c)
	familyId := authUser.FamilyId

	ctx := c.Request().Context()

	uploadBatch, err := con.mediaItemService.CreateUploadBatch(ctx, familyId, albumId)
	if err != nil {
		return err
	}

	return c.JSON(200, models.CreateUploadBatchResponse{
		ID:        uploadBatch.ID,
		AlbumID:   uploadBatch.AlbumID,
		CreatedAt: uploadBatch.CreatedAt,
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
	file, err := c.FormFile("file")
	if err != nil {
		return apperr.NewValidationError("message.validation.required", map[string]string{"field": "message-item.file"})
	}

	albumId, err := utils.ConvertToInt32(c.QueryParam("album_id"))
	if err != nil {
		return apperr.NewValidationError("message.validation.required", map[string]string{"field": "message-item.album_id"})
	}

	uploadBatchID, err := utils.ConvertToInt32(c.QueryParam("upload_batch_id"))
	if err != nil {
		return apperr.NewValidationError("message.validation.required", map[string]string{"field": "message-item.upload_batch_id"})
	}

	retry, _ := utils.ConvertToBool(c.QueryParam("retry"))

	authUser := middlewares.GetAuthUser(c)
	familyId := authUser.FamilyId
	ctx := c.Request().Context()

	result, err := con.mediaItemService.UploadImage(ctx, file, familyId, albumId, uploadBatchID, retry)
	if err != nil {
		log.Println("upload image error: ", err)
		return err
	}

	return c.JSON(200, models.UploadImageResponse{
		MediaItemID: result.MediaItemID,
		Status:      result.Status,
	})
}

// @Tags MediaItem
// @Description 사진이 있는 년도와 월 목록 조회
// @Router /media-item/range [get]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Success 200 {object} []services.MediaItemRange
func (con *MediaItemHandler) GetMediaItemRange(c echo.Context) error {

	authUser := middlewares.GetAuthUser(c)
	ranges, err := con.mediaItemService.GetMediaItemRange(c.Request().Context(), authUser.GroupId)
	if err != nil {
		return err
	}
	return c.JSON(200, ranges)
}

// @Tags MediaItem
// @Description 미디어 아이템 목록 조회
// @Router /media-item [get]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Param from query string true "From" example(2025-01-01)
// @Param to query string true "To" example(2025-01-01)
// @Param identity_ids query []int false "Identity IDs"
// @Param album_id query int false "Album ID"
// @Success 200
func (con *MediaItemHandler) GetMediaItems(c echo.Context) error {
	reqParams := new(models.GetMediaItemsRequest)
	if err := c.Bind(reqParams); err != nil {
		log.Println("bind get media items request error: ", err)
		return err
	}
	if err := c.Validate(reqParams); err != nil {
		return err
	}

	authUser := middlewares.GetAuthUser(c)
	ctx := c.Request().Context()

	mediaItems, err := con.mediaItemService.GetMediaItemsByTakenAt(ctx, authUser.GroupId, *reqParams.From, *reqParams.To, reqParams.AlbumID, reqParams.IdentityIDs)
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
		return apperr.NewValidationError("message.validation.required", map[string]string{"field": "message-item.upload_batch_id"})
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
