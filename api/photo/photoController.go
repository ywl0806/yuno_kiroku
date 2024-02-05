package photo

import (
	"log"

	"github.com/labstack/echo/v4"
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
// @Router /photo [post]
func (con *PhotoController) UploadPhoto(c echo.Context) error {
	file, err := c.FormFile("file")

	if err != nil {
		return err
	}
	_, err = con.standardStorage.SaveFile(file, "")
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
