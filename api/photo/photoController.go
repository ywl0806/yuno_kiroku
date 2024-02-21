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
		log.Println("formfile error : ", err)
		return err
	}
	ext := strings.Split(file.Filename, ".")[1]
	originalFile, err := file.Open()
	if err != nil {
		log.Println("file open error: ", err)
		return err
	}
	defer originalFile.Close()
	// file resize
	// convert to jpeg
	resizedFile := bytes.NewBuffer(nil)
	err = imageHandler.ResizeImage(originalFile, resizedFile, ext)
	if err != nil {
		log.Println("resize error: ", err)

		return err
	}
	originalFilename := strings.Split(file.Filename, ".")[0]

	originalUrl, err := con.standardStorage.SaveFile(resizedFile, "", originalFilename+".jpeg")
	if err != nil {
		log.Println("Standard Storage Error: ", err)
		return err
	}

	thumbnailUrl, err := con.longTermStorage.SaveFile(originalFile, "", file.Filename)
	if err != nil {
		log.Println("Longterm Storage Error: ", err)
		return err
	}

	var photo = models.Photo{
		ThumbnailUrl: thumbnailUrl,
		OriginalUrl:  originalUrl,
		FileName:     file.Filename,
		CreatedBy:    "admin",
		UpdatedBy:    "admin",
	}

	newPhoto, err := con.photoStore.CreatePicture(photo)

	if err != nil {
		log.Println(err)
		return err
	}

	return c.JSON(200, newPhoto)
}
