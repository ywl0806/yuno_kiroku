package photo

import (
	"database/sql"
	"log"
	"mime/multipart"

	"github.com/labstack/echo/v4"

	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
	"github.com/ywl0806/yuno_kiroku/internal/api/utils"
	"github.com/ywl0806/yuno_kiroku/internal/db"
)

type PhotoHandler struct {
	photoService *PhotoService
}

func NewPhotoHandler(
	photoService *PhotoService,
) *PhotoHandler {
	return &PhotoHandler{
		photoService: photoService,
	}
}

type UploadPhotoRequest struct {
	File        *multipart.FileHeader `form:"file" validate:"required"`
	ClanGroupId *int32                `query:"clan_group_id"`
}

func (UploadPhotoRequest) bind(c echo.Context) error {
	var req UploadPhotoRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}
	return nil
}

// @Description 사진 업로드
// @Accept  multipart/form-data
// @Param file formData file true "file"
// @Param clan_group_id query string false "Clan Group ID"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /photo/upload [post]
func (con *PhotoHandler) UploadPhoto(c echo.Context) error {
	req := new(UploadPhotoRequest)

	if err := req.bind(c); err != nil {
		return err
	}
	// 그룹ID 가져오기
	authUser := middlewares.GetAuthUser(c)
	groupId := authUser.GroupId

	// 업로드 경로 생성
	uploadPath := con.photoService.CreateUploadPath(groupId, req.ClanGroupId)

	// 사진 업로드
	uploadResult, err := con.photoService.UploadPhoto(req.File, uploadPath)

	if err != nil {
		log.Println("사진 업로드 실패: ", err)
		return echo.NewHTTPError(400, err.Error())
	}

	// 사진 저장 파라미터 생성
	var params = db.CreatePhotoParams{}
	utils.ConvertStruct(&params, uploadResult)

	if req.ClanGroupId != nil {
		params.ClanGroupID = sql.NullInt32{
			Int32: *req.ClanGroupId,
			Valid: true,
		}
	}
	// 사진 저장
	photo, err := con.photoService.CreatePhoto(c.Request().Context(), params)

	if err != nil {
		log.Println("사진 저장 실패: ", err)
		return err
	}

	return c.JSON(200, photo)
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

// @Description get photo group by date
// @Router /photo/group [get]
// @Param from query string false "from" format(date-time) example(2024-01-01T00:00:00Z)
// @Param to query string false "to" format(date-time) example(2024-05-01T00:00:00Z)
// func (con *PhotoHandler) GetPhotosGroup(c echo.Context) error {
// 	fromQ := c.QueryParam("from")

// 	toQ := c.QueryParam("to")

// 	from := utils.GetDateFromStr(fromQ)
// 	to := utils.GetDateFromStr(toQ)
// 	groupId, err := utils.ContextInt32(c, consts.UserGroupIdKey)

// 	clanGroupId, err := utils.QueryParamInt32(c, "clan_group_id")
// 	if err != nil {
// 		clanGroupId = 0
// 	}

// 	params := db.FindPhotosByPhotoCreatedAtParams{
// 		GroupID:            int32(groupId),
// 		PhotoCreatedAtFrom: from,
// 		PhotoCreatedAtTo:   to,
// 		ClanGroupID: sql.NullInt32{
// 			Int32: clanGroupId,
// 			Valid: true,
// 		},
// 	}

// 	photos, err := con.photoService.FindPhotosByPhotoCreatedAt(c.Request().Context(), &params)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	return c.JSON(200, photos)
// }

// @Description get first photo
// @Router /photo/first [get]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Success 200
// func (con *PhotoHandler) GetFirstPhoto(c echo.Context) error {
// 	opts := options.FindOne()
// 	opts.SetSort(map[string]int{"photo_created_at": 1})

// 	photo, err := con.photoStore.FindOnePhoto(opts)

// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	return c.JSON(200, photo)
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
