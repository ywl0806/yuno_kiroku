package photo

import (
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

func (con *PhotoController) PutPhoto(c echo.Context) error {
	file, err := c.FormFile("file")

	if err != nil {
		return err
	}
	con.standardStorage.SaveFile(file, "")
	return nil
}
