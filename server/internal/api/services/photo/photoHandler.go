package photo

import (
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ywl0806/yuno_kiroku/internal/api/consts"
	"github.com/ywl0806/yuno_kiroku/internal/api/services/photo/models"
	"github.com/ywl0806/yuno_kiroku/internal/api/services/photo/store"
	"github.com/ywl0806/yuno_kiroku/internal/api/utils"
	"github.com/ywl0806/yuno_kiroku/internal/db"
)

type PhotoHandler struct {
	photoStore   *store.PhotoStore
	photoService *PhotoService
}

func NewPhotoHandler(
	photoStore *store.PhotoStore, photoService *PhotoService,
) *PhotoHandler {
	return &PhotoHandler{
		photoStore:   photoStore,
		photoService: photoService,
	}
}

// @Description 사진 업로드
// @Accept  multipart/form-data
// @Param file formData file true "file"
// @Param clan_group_id path string false "Clan Group ID"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /photo/upload [post]
func (con *PhotoHandler) UploadPhoto(c echo.Context) error {
	file, err := c.FormFile("file")

	if err != nil {
		log.Println("file open error: ", err)
		return err
	}
	// 그룹ID 가져오기
	groupId, err := strconv.Atoi(c.Get(consts.UserGroupIdKey).(string))
	if err != nil {
		return c.JSON(400, "user group id is required")
	}

	// 클랜 그룹ID 가져오기
	clanGroupId, err := strconv.Atoi(c.Param("clan_group_id"))

	// 업로드 경로 생성
	uploadPath := con.photoService.CreateUploadPath(int32(groupId), &int32(clanGroupId))

	// 사진 업로드
	uploadResult, err := con.photoService.UploadPhoto(file, uploadPath)

	if err != nil {
		log.Println("사진 업로드 실패: ", err)
		return err
	}

	// 사진 저장 파라미터 생성
	var params = db.CreatePhotoParams{}
	utils.ConvertStruct(&params, uploadResult)

	params.ClanGroupID = int32(clanGroupId)
	// 사진 저장
	photo, err := con.photoService.CreatePhoto(c.Request().Context(), params)

	if err != nil {
		log.Println("사진 저장 실패: ", err)
		return err
	}

	return c.JSON(200, photo)
}

// @Description get photo list
// @Router /photo [get]
// @Param limit query int false "limit"
// @Param skip query int false "skip"
func (con *PhotoHandler) GetPhotoList(c echo.Context) error {

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}
	skip, err := strconv.Atoi(c.QueryParam("skip"))
	if err != nil {
		skip = 0
	}
	parmas := store.FindPictureParams{
		Limit: &limit,
		Skip:  &skip,
	}
	photos, err := con.photoStore.FindPictures(&parmas)
	if err != nil {
		log.Println(err)
		return err
	}
	return c.JSON(200, photos)
}

// @Description get photo group by date
// @Router /photo/group [get]
// @Param from query string false "from" format(date-time) example(2024-01-01T00:00:00Z)
// @Param to query string false "to" format(date-time) example(2024-05-01T00:00:00Z)
func (con *PhotoHandler) GetPhotosGroup(c echo.Context) error {
	fromQ := c.QueryParam("from")

	toQ := c.QueryParam("to")

	from := utils.GetDateFromStr(fromQ)
	to := utils.GetDateFromStr(toQ)
	groupIdInterface := c.Get(consts.UserGroupIdKey)

	var groupId string
	if groupIdInterface != nil {
		groupId = groupIdInterface.(string)
	}

	photos, err := con.photoStore.FindPicturesGroupByDate(from, to, groupId, c.Request().Context())
	if err != nil {
		log.Println(err)
		return err
	}
	return c.JSON(200, photos)
}

// @Description get photo range
// @Router /photo/range [get]
func (con *PhotoHandler) GetPhotoRange(c echo.Context) error {
	photos, err := con.photoStore.FindPhotosRange()
	if err != nil {
		log.Println(err)
		return err
	}

	return c.JSON(200, photos)
}

// @Description get first photo
// @Router /photo/first [get]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Success 200
func (con *PhotoHandler) GetFirstPhoto(c echo.Context) error {
	opts := options.FindOne()
	opts.SetSort(map[string]int{"photo_created_at": 1})

	photo, err := con.photoStore.FindOnePhoto(opts)

	if err != nil {
		log.Println(err)
		return err
	}
	return c.JSON(200, photo)
}

// @Description upload live photo
// @Accept  multipart/form-data
// @Param photo formData file true "photo"
// @Param live formData file true "live"
// @Param clan_group_id path string false "Clan Group ID"
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Router /photo/upload-live [post]
func (con *PhotoHandler) UploadLivePhoto(c echo.Context) error {
	photoFile, err := c.FormFile("photo")
	if err != nil {
		log.Println("no photo formfile error: ", err)
		return c.JSON(400, err)
	}

	liveMovie, err := c.FormFile("live")
	if err != nil {
		log.Println("no live movie formfile error: ", err)
		return c.JSON(400, err)
	}
	groupId := c.Get(consts.UserGroupIdKey)

	if groupId == nil {
		return c.JSON(400, "user group id is required")
	}

	clanGroupId := c.Param("clan_group_id")

	uploadPath := groupId.(string)
	if clanGroupId != "" {
		uploadPath += "/" + clanGroupId
	}

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

	photo := models.Photo{
		ThumbnailUrl:    uploadPhotoResult.ThumbnailUrl,
		OriginalUrl:     uploadPhotoResult.OriginalUrl,
		LiveUrl:         uploadLiveMovieResult.LiveUrl,
		OriginalLiveUrl: uploadLiveMovieResult.OriginalLiveUrl,
		FileName:        photoFile.Filename,
		PhotoCreatedAt:  uploadPhotoResult.PhotoCreatedAt,
		Width:           uploadPhotoResult.Width,
		Height:          uploadPhotoResult.Height,
		Orientation:     uploadPhotoResult.Orientation,
		CreatedBy:       "admin",
		UpdatedBy:       "admin",
	}

	photo, err = con.photoStore.CreatePhoto(photo)
	if err != nil {
		log.Println("upload photo error: ", err)
		return err
	}
	return c.JSON(200, photo)

}
