package photo

import "github.com/labstack/echo/v4"

func Register(root *echo.Group, photoController PhotoController) {
	photoRouter := root.Group("/photo")

	photoRouter.POST("/upload", photoController.UploadPhoto)
	photoRouter.GET("/list", photoController.GetPhotoList)
}
