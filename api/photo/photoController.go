package photo

import (
	"bytes"
	"log"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/ywl0806/yuno_kiroku/api/lib/imageHandler"
	"github.com/ywl0806/yuno_kiroku/api/lib/storage"
	"github.com/ywl0806/yuno_kiroku/api/photo/models"
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

// @Description upload photo
// @Accept  multipart/form-data
// @Param file formData file true "file"
// @Router /photo/upload [post]
func (con *PhotoController) UploadPhoto(c echo.Context) error {
	file, err := c.FormFile("file")

	if err != nil {
		log.Fatal("formfile error : ", err)
		return err
	}
	ext := strings.Split(file.Filename, ".")[1]
	originalFile, _ := file.Open()

	// file resize
	// convert to jpeg
	resizedFile := bytes.NewBuffer(nil)
	imageHandler.ResizeImage(originalFile, resizedFile, ext)

	originalFilename := strings.Split(file.Filename, ".")[0]

	originalUrl, err := con.standardStorage.SaveFile(resizedFile, "", originalFilename+".jpeg")
	if err != nil {
		log.Fatal(err)
	}

	thumbnailUrl, err := con.longTermStorage.SaveFile(originalFile, "", file.Filename)
	if err != nil {
		log.Fatal(err)
	}

	var photo = models.Photo{
		ThumbnailUrl: thumbnailUrl,
		OriginalUrl:  originalUrl,
		FileName:     file.Filename,
		CreatedBy:    "admin",
		UpdatedBy:    "admin",
	}

	con.photoStore.CreatePicture(photo)

	return nil
}
