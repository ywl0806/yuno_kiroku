package handlers

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/ywl0806/yuno_kiroku/internal/api/consts"
	customErrors "github.com/ywl0806/yuno_kiroku/internal/api/errors"
	"github.com/ywl0806/yuno_kiroku/internal/api/handlers/models"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
	"github.com/ywl0806/yuno_kiroku/internal/api/services"
	"github.com/ywl0806/yuno_kiroku/internal/api/utils"
	"github.com/ywl0806/yuno_kiroku/internal/db"
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
// @Description 이미지 업로드
// @Accept  multipart/form-data
// @Param file formData file true "file"
// @Param album_id query string true "Album ID"
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

	retry, _ := utils.ConvertToBool(c.QueryParam("retry"))

	// 그룹 ID 가져오기
	authUser := middlewares.GetAuthUser(c)
	groupId := authUser.GroupId

	// 컨텍스트 가져오기
	ctx := c.Request().Context()

	// 이미지 핸들러 생성(이미지 리사이즈)
	imgHandler, err := con.mediaItemService.HandleImage(file)
	if err != nil {
		return err
	}
	// 얼굴 인식 결과 가져오기(이미지 리사이즈 된 이미지 이용)
	faceDetections, err := con.faceService.GetFaceDetection(imgHandler.ResizedFile)
	if err != nil {
		return err
	}

	// 사진 중복체크(재업로드가 아닌 경우만 체크)
	if !retry {
		err = con.faceService.CheckImageDuplicateByFaceDetection(ctx, groupId, albumId, *faceDetections)
		if err != nil {
			return err
		}
	}

	// 업로드 경로 생성
	uploadPath := con.mediaItemService.CreateUploadPath(groupId, albumId)

	// 썸네일 이미지와 원본 이미지 업로드
	thumbnailStorageKey, originalStorageKey, err := con.mediaItemService.UploadImage(imgHandler, uploadPath)
	if err != nil {
		return err
	}

	// 사진 저장 파라미터 생성
	createMediaItemParams := services.CreateMediaItemParams{
		ThumbnailStorageKey: thumbnailStorageKey,
		OriginalStorageKey:  originalStorageKey,
		ImageHandler:        imgHandler,
		GroupId:             groupId,
		AlbumId:             albumId,
		OriginalFilename:    file.Filename,
	}

	// 사진 저장
	mediaItem, err := con.mediaItemService.CreateMediaItem(ctx, &createMediaItemParams)
	if err != nil {
		return err
	}
	// 얼굴 인식 결과 저장
	_, err = con.faceService.SearchAndSaveFaceDetections(ctx, groupId, mediaItem.ID, *faceDetections)
	if err != nil {
		log.Println("얼굴 인식 결과 저장 실패: ", err)
		// 얼굴 인식 실패해도 사진은 저장되었으므로 계속 진행
	}

	return c.JSON(200, "ok")
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
