package photo

import "github.com/labstack/echo/v4"

func Register(root *echo.Group, photoController PhotoController) {
	photoRouter := root.Group("/photo")

	photoRouter.POST("/upload", photoController.UploadPhoto)
	photoRouter.GET("", photoController.GetPhotoList)
	photoRouter.GET("/group", photoController.GetPhotosGroup)
	photoRouter.GET("/range", photoController.GetPhotoRange)
	photoRouter.GET("/first", photoController.GetFirstPhoto)
}
