package handlers

import (
	"log"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/ywl0806/yuno_kiroku/internal/api/handlers/models"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
	"github.com/ywl0806/yuno_kiroku/internal/services"

	"github.com/ywl0806/yuno_kiroku/internal/apperr"
	"github.com/ywl0806/yuno_kiroku/pkg/utils"
)

type MediaItemHandler struct {
	mediaItemService *services.MediaItemService
}

func NewMediaItemHandler(
	mediaItemService *services.MediaItemService,
) *MediaItemHandler {
	return &MediaItemHandler{
		mediaItemService: mediaItemService,
	}
}

// @Tags MediaItem
// @Description S3 직접 업로드를 위한 Presigned PUT URL 발급
// @Accept json
// @Param body body models.PresignedUploadRequest true "파일 정보"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /media-item/presigned-url [post]
func (con *MediaItemHandler) CreatePresignedUpload(c echo.Context) error {
	req := new(models.PresignedUploadRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	authUser := middlewares.GetAuthUser(c)
	ctx := c.Request().Context()

	result, err := con.mediaItemService.CreatePresignedUpload(ctx, req.FileName, req.ContentType, authUser.FamilyId, req.AlbumID, req.UploadBatchID)
	if err != nil {
		log.Println("create presigned upload error:", err)
		return err
	}

	return c.JSON(200, models.PresignedUploadResponse{
		MediaItemID:  result.MediaItemID,
		PresignedURL: result.PresignedURL,
		StorageKey:   result.StorageKey,
		ExpiresIn:    3600,
	})
}

// @Tags MediaItem
// @Description S3 직접 업로드를 위한 Presigned PUT URL 배치 발급
// @Accept json
// @Param body body models.BatchPresignedUploadRequest true "파일 정보 배열"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /media-item/presigned-urls [post]
func (con *MediaItemHandler) CreateBatchPresignedUpload(c echo.Context) error {
	req := new(models.BatchPresignedUploadRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	authUser := middlewares.GetAuthUser(c)
	ctx := c.Request().Context()

	items := make([]services.BatchPresignedUploadItem, len(req.Files))
	for i, f := range req.Files {
		items[i] = services.BatchPresignedUploadItem{
			FileName:    f.FileName,
			ContentType: f.ContentType,
		}
	}

	results, err := con.mediaItemService.CreateBatchPresignedUpload(ctx, items, authUser.FamilyId, req.AlbumID, req.UploadBatchID)
	if err != nil {
		log.Println("create batch presigned upload error:", err)
		return err
	}

	successItems := make([]models.BatchPresignedUploadSuccessItem, len(results.Success))
	for i, r := range results.Success {
		successItems[i] = models.BatchPresignedUploadSuccessItem{
			MediaItemID:   r.MediaItemID,
			PresignedURL:  r.PresignedURL,
			StorageKey:    r.StorageKey,
			ExpiresIn:     3600,
			OriginalIndex: r.OriginalIndex,
		}
	}

	failedItems := make([]models.BatchPresignedUploadFailedItem, len(results.Failed))
	for i, f := range results.Failed {
		failedItems[i] = models.BatchPresignedUploadFailedItem{FileName: f.FileName, Index: f.Index, Reason: f.Reason}
	}

	return c.JSON(200, models.BatchPresignedUploadResponse{Success: successItems, Failed: failedItems})
}

// @Tags MediaItem
// @Description 업로드 배치 생성
// @Param album_id query string true "Album ID"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /media-item/upload-batch [post]
func (con *MediaItemHandler) CreateUploadBatch(c echo.Context) error {
	albumId := c.QueryParam("album_id")
	if albumId == "" {
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
// @Param album_id query string false "Album ID"
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

	mediaItems, err := con.mediaItemService.GetMediaItemsByTakenAt(ctx, authUser.GroupId, authUser.ID, *reqParams.From, *reqParams.To)
	if err != nil {
		log.Println("get media items by taken at error: ", err)
		return err
	}
	return c.JSON(200, models.NewMediaItemsResponse(mediaItems))
}

// @Tags MediaItem
// @Description 검색 (페이지네이션)
// @Router /media-item/search [get]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Param from query string false "From (RFC3339)"
// @Param to query string false "To (RFC3339)"
// @Param identity_ids query []int false "Identity IDs"
// @Param album_id query string false "Album ID"
// @Param page query int false "Page (1-based, default 1)"
// @Success 200 {object} models.SearchMediaItemsResponse
func (con *MediaItemHandler) SearchMediaItems(c echo.Context) error {
	reqParams := new(models.SearchMediaItemsRequest)
	if err := c.Bind(reqParams); err != nil {
		return err
	}

	authUser := middlewares.GetAuthUser(c)
	ctx := c.Request().Context()

	result, err := con.mediaItemService.SearchMediaItems(ctx, authUser.GroupId, authUser.ID, reqParams.From, reqParams.To, reqParams.AlbumID, reqParams.IdentityIDs, reqParams.Liked, reqParams.TagIDs, reqParams.Page)
	if err != nil {
		log.Println("search media items error:", err)
		return err
	}

	items := make([]models.MediaItemResponse, len(result.Items))
	for i, item := range result.Items {
		items[i] = *models.NewSearchMediaItemResponse(&item)
	}
	return c.JSON(200, models.SearchMediaItemsResponse{
		Items:   items,
		HasNext: result.HasNext,
		Page:    reqParams.Page,
	})
}

// @Tags MediaItem
// @Description 업로드 배치 목록 + 썸네일 5개 조회
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /media-item/upload-batch [get]
func (con *MediaItemHandler) GetUploadBatches(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)
	ctx := c.Request().Context()
	page := 1
	if p := c.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	batches, hasNext, err := con.mediaItemService.GetUploadBatchesWithThumbnails(ctx, authUser.GroupId, page)
	if err != nil {
		return err
	}

	items := make([]models.UploadBatchWithThumbnailsResponse, len(batches))
	for i, b := range batches {
		items[i] = models.NewUploadBatchWithThumbnailsResponse(b)
	}
	return c.JSON(200, models.GetUploadBatchesResponse{
		Items:   items,
		HasNext: hasNext,
		Page:    page,
	})
}

// @Tags MediaItem
// @Description 특정 배치의 미디어 아이템 목록 조회 (페이지네이션)
// @Param id path int true "Upload Batch ID"
// @Param page query int false "Page (1-based, default 1)"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /media-item/upload-batch/{id}/items [get]
func (con *MediaItemHandler) GetUploadBatchItems(c echo.Context) error {
	uploadBatchID, err := utils.ConvertToInt32(c.Param("id"))
	if err != nil {
		return apperr.NewValidationError("message.validation.required", map[string]string{"field": "id"})
	}

	page := 1
	if p := c.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	authUser := middlewares.GetAuthUser(c)
	ctx := c.Request().Context()

	result, err := con.mediaItemService.GetMediaItemsByUploadBatch(ctx, uploadBatchID, authUser.ID, page)
	if err != nil {
		log.Println("get upload batch items error:", err)
		return err
	}

	items := make([]models.MediaItemResponse, len(result.Items))
	for i, item := range result.Items {
		items[i] = *models.NewUploadBatchItemResponse(&item)
	}
	return c.JSON(200, models.SearchMediaItemsResponse{
		Items:   items,
		HasNext: result.HasNext,
		Page:    page,
	})
}

// @Tags MediaItem
// @Description 미디어 아이템 삭제
// @Param id path string true "MediaItem ID"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /media-item/{id} [delete]
func (con *MediaItemHandler) DeleteMediaItem(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return apperr.NewValidationError("message.validation.required", map[string]string{"field": "id"})
	}
	authUser := middlewares.GetAuthUser(c)
	ctx := c.Request().Context()
	if err := con.mediaItemService.DeleteMediaItem(ctx, id, authUser.FamilyId, authUser.ID); err != nil {
		return err
	}
	return c.NoContent(204)
}

// @Tags MediaItem
// @Description 미디어 아이템 앨범 변경
// @Param id path string true "MediaItem ID"
// @Param body body models.UpdateMediaItemAlbumRequest true "앨범 변경 요청"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /media-item/{id}/album [patch]
func (con *MediaItemHandler) UpdateMediaItemAlbum(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return apperr.NewValidationError("message.validation.required", map[string]string{"field": "id"})
	}
	req := new(models.UpdateMediaItemAlbumRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}
	authUser := middlewares.GetAuthUser(c)
	ctx := c.Request().Context()
	if err := con.mediaItemService.UpdateMediaItemAlbum(ctx, id, authUser.FamilyId, authUser.ID, req.AlbumID); err != nil {
		return err
	}
	return c.NoContent(204)
}

// @Tags MediaItem
// @Description 업로드 배치 상태 조회
// @Param upload_batch_id query string true "Upload Batch ID"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /media-item/upload-batch/status [get]
func (con *MediaItemHandler) GetUploadBatchStatus(c echo.Context) error {
	uploadBatchID, err := utils.ConvertToInt32(c.QueryParam("upload_batch_id"))
	if err != nil {
		return apperr.NewValidationError("message.validation.required", map[string]string{"field": "message-item.upload_batch_id"})
	}

	ctx := c.Request().Context()

	uploadStatuses, err := con.mediaItemService.GetUploadStatuses(ctx, uploadBatchID)
	if err != nil {
		return err
	}

	uploadStatusesResponse := models.NewUploadBatchStatusResponse(uploadStatuses)

	return c.JSON(200, uploadStatusesResponse)
}
