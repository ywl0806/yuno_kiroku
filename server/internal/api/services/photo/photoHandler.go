package photo

import (
	"database/sql"
	"log"

	"github.com/labstack/echo/v4"

	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
	"github.com/ywl0806/yuno_kiroku/internal/api/services/photo/models"
	"github.com/ywl0806/yuno_kiroku/internal/api/utils"
	"github.com/ywl0806/yuno_kiroku/internal/db"
)

type PhotoHandler struct {
	photoService *PhotoService
	faceService  *FaceService
}

func NewPhotoHandler(
	photoService *PhotoService,
	faceService *FaceService,
) *PhotoHandler {
	return &PhotoHandler{
		photoService: photoService,
		faceService:  faceService,
	}
}

// @Description 사진 업로드
// @Accept  multipart/form-data
// @Param file formData file true "file"
// @Param clan_group_id query string false "Clan Group ID"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /photo/upload [post]
func (con *PhotoHandler) UploadPhoto(c echo.Context) error {
	// 파일 가져오기
	file, err := c.FormFile("file")
	if err != nil {
		log.Println("파일 가져오기 실패: ", err)
		return echo.NewHTTPError(400, "file is required")
	}

	// clan_group_id 쿼리 파라미터 가져오기
	var clanGroupId *int32
	if clanGroupIdStr := c.QueryParam("clan_group_id"); clanGroupIdStr != "" {
		clanGroupIdInt, err := utils.ConvertToInt32(clanGroupIdStr)
		if err == nil {
			clanGroupId = &clanGroupIdInt
		}
	}

	authUser := middlewares.GetAuthUser(c)
	groupId := authUser.GroupId

	// 업로드 경로 생성
	uploadPath := con.photoService.CreateUploadPath(groupId, clanGroupId)
	// 사진 업로드
	uploadResult, err := con.photoService.UploadPhoto(file, uploadPath)

	if err != nil {
		log.Println("사진 업로드 실패: ", err)
		return echo.NewHTTPError(400, err.Error())
	}

	// 사진 저장 파라미터
	var params = db.CreatePhotoParams{}

	params.FileName = uploadResult.FileName
	params.PhotoCreatedAt = uploadResult.PhotoCreatedAt
	params.ThumbnailUrl = uploadResult.ThumbnailUrl
	params.Width = uploadResult.Width
	params.Height = uploadResult.Height
	params.Orientation = uploadResult.Orientation
	params.GroupID = groupId

	if uploadResult.OriginalUrl != "" {
		params.OriginalUrl = sql.NullString{String: uploadResult.OriginalUrl, Valid: true}
	}

	if clanGroupId != nil {
		params.ClanGroupID = sql.NullInt32{
			Int32: *clanGroupId,
			Valid: true,
		}
	}

	// 사진 저장
	photo, err := con.photoService.CreatePhoto(c.Request().Context(), params)

	if err != nil {
		log.Println("사진 저장 실패: ", err)
		return err
	}

	// 얼굴 인식 결과를 저장
	_, err = con.faceService.SearchAndSaveFaceDetections(c.Request().Context(), groupId, photo.ID, uploadResult.FaceDetections)
	if err != nil {
		log.Println("얼굴 인식 결과 저장 실패: ", err)
		return err
	}

	faceDetectionRows, err := con.faceService.GetFaceDetections(c.Request().Context(), photo.ID)
	if err != nil {
		return err
	}

	uploadPhotoResponse := models.NewUploadPhotoResponse(photo, faceDetectionRows)
	return c.JSON(200, uploadPhotoResponse)
}

// @Description 사진이 있는 년도와 월 목록 조회
// @Router /photo/range [get]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Success 200 {object} []db.GetPhotoRangeRow
func (con *PhotoHandler) GetPhotoRange(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)
	ranges, err := con.photoService.GetPhotoRange(c.Request().Context(), authUser.GroupId, authUser.ClanGroupId)
	if err != nil {
		log.Println("사진 년도와 월 목록 조회 실패: ", err)
		return err
	}
	return c.JSON(200, ranges)
}

// type GetPhotosGroupRequest struct {
// 	From *time.Time `query:"from" validate:"required"`
// 	To   *time.Time `query:"to" validate:"required"`
// }

// func (GetPhotosGroupRequest) bind(c echo.Context, p *db.FindPhotosByPhotoCreatedAtParams) error {
// 	reqParams := new(GetPhotosGroupRequest)
// 	if err := c.Bind(reqParams); err != nil {
// 		return err
// 	}
// 	if err := c.Validate(reqParams); err != nil {
// 		return err
// 	}

// 	authUser := middlewares.GetAuthUser(c)

// 	p.GroupID = authUser.GroupId
// 	p.PhotoCreatedAtFrom = *reqParams.From
// 	p.PhotoCreatedAtTo = *reqParams.To

// 	if authUser.ClanGroupId != nil {
// 		p.ClanGroupID = *authUser.ClanGroupId
// 	}

// 	return nil
// }

// @Description upload live photo
// @Accept  multipart/form-data
// @Param photo formData file true "photo"
// @Param live formData file true "live"
// @Param clan_group_id path string false "Clan Group ID"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /photo/upload-live [post]
// func (con *PhotoHandler) UploadLivePhoto(c echo.Context) error {
// 	photoFile, err := c.FormFile("photo")
// 	if err != nil {
// 		log.Println("no photo formfile error: ", err)
// 		return c.JSON(400, err)
// 	}

// 	liveMovie, err := c.FormFile("live")
// 	if err != nil {
// 		log.Println("no live movie formfile error: ", err)
// 		return c.JSON(400, err)
// 	}
// 	groupId := c.Get(consts.UserGroupIdKey)

// 	if groupId == nil {
// 		return c.JSON(400, "user group id is required")
// 	}

// 	clanGroupId := c.Param("clan_group_id")

// 	uploadPath := groupId.(string)
// 	if clanGroupId != "" {
// 		uploadPath += "/" + clanGroupId
// 	}

// 	uploadPhotoResult, err := con.photoService.UploadPhoto(photoFile, uploadPath)
// 	if err != nil {
// 		log.Println("upload photo error: ", err)
// 		return err
// 	}
// 	uploadLiveMovieResult, err := con.photoService.UploadLiveMovie(liveMovie)

// 	if err != nil {
// 		log.Println("upload live photo error: ", err)
// 		return err
// 	}

// 	photo := models.Photo{
// 		ThumbnailUrl:    uploadPhotoResult.ThumbnailUrl,
// 		OriginalUrl:     uploadPhotoResult.OriginalUrl,
// 		LiveUrl:         uploadLiveMovieResult.LiveUrl,
// 		OriginalLiveUrl: uploadLiveMovieResult.OriginalLiveUrl,
// 		FileName:        photoFile.Filename,
// 		PhotoCreatedAt:  uploadPhotoResult.PhotoCreatedAt,
// 		// Width:           uploadPhotoResult.Width,
// 		// Height:          uploadPhotoResult.Height,
// 		// Orientation:     uploadPhotoResult.Orientation,
// 		CreatedBy: "admin",
// 		UpdatedBy: "admin",
// 	}

// 	photo, err = con.photoStore.CreatePhoto(photo)
// 	if err != nil {
// 		log.Println("upload photo error: ", err)
// 		return err
// 	}
// 	return c.JSON(200, photo)

// }
