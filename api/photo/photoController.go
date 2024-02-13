package photo

import (
	"bytes"
	"log"

	"github.com/labstack/echo/v4"

	"github.com/ywl0806/yuno_kiroku/api/lib/imageHandler"
	"github.com/ywl0806/yuno_kiroku/api/lib/storage"
	"github.com/ywl0806/yuno_kiroku/api/photo/store"
)

type PhotoController struct {
	photoStore      store.PhotoStore
	standardStorage storage.StorageService
	longTermStorage storage.StorageService
}

func NewPhotoController(
	photoStore store.PhotoStore, standardStorage storage.StorageService, longTermStorage storage.StorageService,
) *PhotoController {
	return &PhotoController{
		photoStore:      photoStore,
		standardStorage: standardStorage,
		longTermStorage: longTermStorage,
	}
}

// @Description put photo
// @Accept  multipart/form-data
// @Param file formData file true "file"
// @Router /photo/upload [post]
func (con *PhotoController) UploadPhoto(c echo.Context) error {
	file, err := c.FormFile("file")

	if err != nil {
		log.Fatal("formfile error : ", err)
		return err
	}

	f, _ := file.Open()

	rf := bytes.NewBuffer(nil)

	imageHandler.ResizeImage(f, rf)

	_, err = con.longTermStorage.SaveFile(rf, "", "test.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	_, err = con.standardStorage.SaveFile(f, "", file.Filename)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
