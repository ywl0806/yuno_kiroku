package photo

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
)

func Register(root *echo.Group, photoController PhotoHandler) {
	photoRouter := root.Group("/photo")
	guard := middlewares.NewGuard()

	photoRouter.POST("/upload/{clan_group_id}", photoController.UploadPhoto, guard.Handler)
	photoRouter.POST("/upload", photoController.UploadPhoto, guard.Handler)
	photoRouter.POST("/upload-live/{clan_group_id}", photoController.UploadLivePhoto, guard.Handler)
	photoRouter.POST("/upload-live", photoController.UploadLivePhoto, guard.Handler)
	photoRouter.GET("", photoController.GetPhotoList)
	photoRouter.GET("/group", photoController.GetPhotosGroup, guard.Handler)
	photoRouter.GET("/range", photoController.GetPhotoRange)
	photoRouter.GET("/first", photoController.GetFirstPhoto, guard.Handler)

}
