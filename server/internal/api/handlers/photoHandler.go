package handlers

import (
	"database/sql"
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

type PhotoHandler struct {
	photoService *services.PhotoService
	faceService  *services.FaceService
}

func NewPhotoHandler(
	photoService *services.PhotoService,
	faceService *services.FaceService,
) *PhotoHandler {
	return &PhotoHandler{
		photoService: photoService,
		faceService:  faceService,
	}
}

// @Tags Photo
// @Description 사진 업로드
// @Accept  multipart/form-data
// @Param file formData file true "file"
// @Param album_id query string true "Album ID"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /photo/upload [post]
func (con *PhotoHandler) UploadPhoto(c echo.Context) error {
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

	authUser := middlewares.GetAuthUser(c)
	groupId := authUser.GroupId

	// 사진 업로드 및 저장 (얼굴 인식 포함)
	result, err := con.photoService.UploadAndSavePhoto(c.Request().Context(), file, groupId, albumId)
	if err != nil {
		return err
	}

	uploadPhotoResponse := models.NewUploadPhotoResponse(result.Photo, result.FaceDetections)
	return c.JSON(200, uploadPhotoResponse)
}

// @Tags Photo
// @Description 사진이 있는 년도와 월 목록 조회
// @Router /photo/range [get]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Success 200 {object} []db.GetPhotoRangeRow
func (con *PhotoHandler) GetPhotoRange(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)
	ranges, err := con.photoService.GetPhotoRange(c.Request().Context(), authUser.ClanGroupId)
	if err != nil {
		log.Println("사진 년도와 월 목록 조회 실패: ", err)
		return err
	}
	return c.JSON(200, ranges)
}

type GetPhotosRequest struct {
	From *time.Time `query:"from" validate:"required"`
	To   *time.Time `query:"to" validate:"required"`
}

func (GetPhotosRequest) bind(c echo.Context, p *db.FindPhotosByPhotoCreatedAtParams) error {
	reqParams := new(GetPhotosRequest)
	if err := c.Bind(reqParams); err != nil {
		return err
	}
	if err := c.Validate(reqParams); err != nil {
		return err
	}

	authUser := middlewares.GetAuthUser(c)

	p.PhotoCreatedAtFrom = *reqParams.From
	p.PhotoCreatedAtTo = *reqParams.To
	p.ClanGroupID = authUser.ClanGroupId
	return nil
}

// @Tags Photo
// @Description 사진 목록 조회
// @Router /photo [get]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Param from query string true "From" example(2025-01-01)
// @Param to query string true "To" example(2025-01-01)
// @Success 200
func (con *PhotoHandler) GetPhotos(c echo.Context) error {
	reqParams := new(GetPhotosRequest)

	params := new(db.FindPhotosByPhotoCreatedAtParams)
	if err := reqParams.bind(c, params); err != nil {
		return err
	}
	photos, err := con.photoService.GetPhotos(c.Request().Context(), params)
	if err != nil {
		return err
	}

	photosResponse := models.NewPhotosResponse(photos)
	return c.JSON(200, photosResponse)

}

// @Tags Photo
// @Description upload live photo
// @Accept  multipart/form-data
// @Param photo formData file true "photo"
// @Param live formData file true "live"
// @Param album_id query string true "Album ID"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /photo/upload-live [post]
func (con *PhotoHandler) UploadLivePhoto(c echo.Context) error {
	photoFile, err := c.FormFile("photo")
	if err != nil {
		log.Println("no photo formfile error: ", err)
		return customErrors.NewRequiredError(consts.Photo)
	}

	liveMovie, err := c.FormFile("live")
	if err != nil {
		log.Println("no live movie formfile error: ", err)
		return customErrors.NewRequiredError(consts.LiveMovie)
	}
	groupId := middlewares.GetAuthUser(c).GroupId

	albumId, err := utils.ConvertToInt32(c.QueryParam("album_id"))
	if err != nil {
		return customErrors.NewRequiredError(consts.Album)
	}
	uploadPath := con.photoService.CreateUploadPath(groupId, albumId)

	uploadPhotoResult, err := con.photoService.UploadPhoto(photoFile, uploadPath)
	if err != nil {
		log.Println("upload photo error: ", err)
		return err
	}
	uploadLiveMovieResult, err := con.photoService.UploadLiveMovie(liveMovie)

	if err != nil {
		log.Println("upload live photo error: ", err)
		return err
	}

	params := db.CreatePhotoParams{
		GroupID:         groupId,
		AlbumID:         albumId,
		ThumbnailUrl:    uploadPhotoResult.ThumbnailUrl,
		OriginalUrl:     sql.NullString{String: uploadPhotoResult.OriginalUrl, Valid: true},
		LiveUrl:         sql.NullString{String: uploadLiveMovieResult.LiveUrl, Valid: true},
		OriginalLiveUrl: sql.NullString{String: uploadLiveMovieResult.OriginalLiveUrl, Valid: true},
		FileName:        photoFile.Filename,
		PhotoCreatedAt:  uploadPhotoResult.PhotoCreatedAt,
		Width:           uploadPhotoResult.Width,
		Height:          uploadPhotoResult.Height,
		Orientation:     uploadPhotoResult.Orientation,
	}

	photo, err := con.photoService.CreatePhoto(c.Request().Context(), params)
	if err != nil {
		log.Println("upload photo error: ", err)
		return err
	}

	return c.JSON(200, models.NewPhotoResponse(photo))

}

// @Tags Photo
// @Description 신원 ID를 기반으로 랜덤 사진 조회
// @Router /photo/identity/{identity_id}/random [get]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Param identity_id path int true "Identity ID"
// @Success 200 {object} models.IdentityRandomPhotoResponse
func (con *PhotoHandler) GetIdentityRandomPhoto(c echo.Context) error {
	identityId, err := utils.ConvertToInt32(c.Param("identity_id"))
	if err != nil {
		return customErrors.NewRequiredError(consts.Identity)
	}
	authUser := middlewares.GetAuthUser(c)
	photo, err := con.photoService.GetIdentityRandomPhoto(c.Request().Context(), authUser.ClanGroupId, identityId)
	if err != nil {
		return err
	}
	return c.JSON(200, models.NewIdentityRandomPhotoResponse(photo))
}
