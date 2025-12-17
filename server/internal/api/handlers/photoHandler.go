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

	// 그룹 ID 가져오기
	authUser := middlewares.GetAuthUser(c)
	groupId := authUser.GroupId

	// 컨텍스트 가져오기
	ctx := c.Request().Context()

	// 이미지 핸들러 생성(이미지 리사이즈)
	imgHandler, err := con.photoService.HandleImage(file)
	if err != nil {
		return err
	}
	// 얼굴 인식 결과 가져오기(이미지 리사이즈 된 이미지 이용)
	faceDetections, err := con.faceService.GetFaceDetection(imgHandler.ResizedFile)
	if err != nil {
		return err
	}

	// 사진 중복체크
	err = con.photoService.CheckPhotoDuplicateByFaceDetection(ctx, groupId, albumId, *faceDetections)
	if err != nil {
		return err
	}

	// 업로드 경로 생성
	uploadPath := con.photoService.CreateUploadPath(groupId, albumId)

	// 썸네일 이미지와 원본 이미지 업로드
	thumbnailUrl, originalUrl, err := con.photoService.UploadPhoto(imgHandler, uploadPath)
	if err != nil {
		return err
	}

	// 사진 저장 파라미터 생성
	createPhotoParams := services.CreatePhotoParams{
		ThumbnailUrl:     thumbnailUrl,
		OriginalUrl:      originalUrl,
		ImageHandler:     imgHandler,
		GroupId:          groupId,
		AlbumId:          albumId,
		OriginalFilename: file.Filename,
	}

	// 사진 저장
	photo, err := con.photoService.CreatePhoto(ctx, &createPhotoParams)
	if err != nil {
		return err
	}
	// 얼굴 인식 결과 저장
	_, err = con.faceService.SearchAndSaveFaceDetections(ctx, groupId, photo.ID, *faceDetections)
	if err != nil {
		log.Println("얼굴 인식 결과 저장 실패: ", err)
		// 얼굴 인식 실패해도 사진은 저장되었으므로 계속 진행
	}

	// 얼굴 인식 결과 조회
	faceDetectionRows, err := con.faceService.GetFaceDetections(ctx, photo.ID)
	if err != nil {
		log.Println("얼굴 인식 결과 조회 실패: ", err)
		// 조회 실패해도 빈 배열로 반환
		faceDetectionRows = []db.GetFaceDetectionsByPhotoIdRow{}
	}

	// 사진 응답 생성
	uploadPhotoResponse := models.NewUploadPhotoResponse(photo, faceDetectionRows)
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
	// 컨텍스트 가져오기
	ctx := c.Request().Context()

	// 이미지 핸들러 생성(이미지 리사이즈)
	imgHandler, err := con.photoService.HandleImage(photoFile)
	if err != nil {
		return err
	}
	// 얼굴 인식 결과 가져오기(이미지 리사이즈 된 이미지 이용)
	faceDetections, err := con.faceService.GetFaceDetection(imgHandler.ResizedFile)
	if err != nil {
		return err
	}

	// 업로드 경로 생성
	uploadPath := con.photoService.CreateUploadPath(groupId, albumId)

	// 썸네일 이미지와 원본 이미지 업로드
	thumbnailUrl, originalUrl, err := con.photoService.UploadPhoto(imgHandler, uploadPath)
	if err != nil {
		return err
	}

	// 라이브 무비 업로드
	liveUrl, originalLiveUrl, err := con.photoService.UploadLiveMovie(liveMovie)
	if err != nil {
		log.Println("upload live movie error: ", err)
		return err
	}

	// 사진 저장 파라미터 생성
	createPhotoParams := services.CreatePhotoParams{
		ThumbnailUrl:     thumbnailUrl,
		OriginalUrl:      originalUrl,
		LiveUrl:          liveUrl,
		OriginalLiveUrl:  originalLiveUrl,
		ImageHandler:     imgHandler,
		GroupId:          groupId,
		AlbumId:          albumId,
		OriginalFilename: photoFile.Filename,
	}

	// 사진 저장
	photo, err := con.photoService.CreatePhoto(ctx, &createPhotoParams)
	if err != nil {
		return err
	}

	// 얼굴 인식 결과 저장
	_, err = con.faceService.SearchAndSaveFaceDetections(ctx, groupId, photo.ID, *faceDetections)
	if err != nil {
		log.Println("얼굴 인식 결과 저장 실패: ", err)
		// 얼굴 인식 실패해도 사진은 저장되었으므로 계속 진행
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
