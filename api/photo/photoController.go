package photo

import (
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ywl0806/yuno_kiroku/api/lib/storage"
	"github.com/ywl0806/yuno_kiroku/api/photo/models"
	"github.com/ywl0806/yuno_kiroku/api/photo/store"
	"github.com/ywl0806/yuno_kiroku/api/utils"
)

type PhotoController struct {
	photoStore      *store.PhotoStore
	standardStorage storage.StorageService
	longTermStorage storage.StorageService
}

func NewPhotoController(
	photoStore *store.PhotoStore, standardStorage storage.StorageService, longTermStorage storage.StorageService,
) *PhotoController {
	return &PhotoController{
		photoStore:      photoStore,
		standardStorage: standardStorage,
		longTermStorage: longTermStorage,
	}
}

// @Description upload photo
// @Accept  multipart/form-data
// @Param file formData file true "file"
// @Router /photo/upload [post]
func (con *PhotoController) UploadPhoto(c echo.Context) error {
	file, err := c.FormFile("file")

	if err != nil {
		log.Println("file open error: ", err)
		return err
	}

	uploadResult, err := con.uploadPhoto(file)

	if err != nil {
		log.Println("upload photo error: ", err)
		return err
	}

	var newPhoto = models.Photo{
		ThumbnailUrl:   uploadResult.ThumbnailUrl,
		OriginalUrl:    uploadResult.OriginalUrl,
		FileName:       uploadResult.FileName,
		PhotoCreatedAt: uploadResult.PhotoCreatedAt,
		Width:          uploadResult.Width,
		Height:         uploadResult.Height,
		CreatedBy:      "admin",
		UpdatedBy:      "admin",
	}

	photo, err := con.photoStore.CreatePhoto(newPhoto)
	if err != nil {
		log.Println("upload photo error: ", err)
		return err
	}
	return c.JSON(200, photo)
}

// @Description get photo list
// @Router /photo [get]
// @Param limit query int false "limit"
// @Param skip query int false "skip"
func (con *PhotoController) GetPhotoList(c echo.Context) error {

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
func (con *PhotoController) GetPhotosGroup(c echo.Context) error {
	fromQ := c.QueryParam("from")

	toQ := c.QueryParam("to")

	from := utils.GetDateFromStr(fromQ)
	to := utils.GetDateFromStr(toQ)

	photos, err := con.photoStore.FindPicturesGroupByDate(from, to)
	if err != nil {
		log.Println(err)
		return err
	}
	return c.JSON(200, photos)
}

// @Description get photo range
// @Router /photo/range [get]
func (con *PhotoController) GetPhotoRange(c echo.Context) error {
	photos, err := con.photoStore.FindPhotosRange()
	if err != nil {
		log.Println(err)
		return err
	}

	return c.JSON(200, photos)
}

// @Description get first photo
// @Router /photo/first [get]
// @Success 200
func (con *PhotoController) GetFirstPhoto(c echo.Context) error {
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
// @Router /photo/upload-live [post]
func (con *PhotoController) UploadLivePhoto(c echo.Context) error {
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

	uploadPhotoResult, err := con.uploadPhoto(photoFile)
	if err != nil {
		log.Println("upload photo error: ", err)
		return err
	}
	uploadLiveMovieResult, err := con.uploadLiveMovie(liveMovie)

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
