package photo

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
)

func Register(root *echo.Group, photoHandler PhotoHandler) {
	photoRouter := root.Group("/photo")
	guard := middlewares.NewGuard()

	photoRouter.POST("/upload/{clan_group_id}", photoHandler.UploadPhoto, guard.Handler)
	photoRouter.POST("/upload", photoHandler.UploadPhoto, guard.Handler)
	// photoRouter.POST("/upload-live/{clan_group_id}", photoHandler.UploadLivePhoto, guard.Handler)
	// photoRouter.POST("/upload-live", photoHandler.UploadLivePhoto, guard.Handler)
	// photoRouter.GET("", photoHandler.GetPhotoList)
	// photoRouter.GET("/group", photoHandler.GetPhotosGroup, guard.Handler)
	photoRouter.GET("/range", photoHandler.GetPhotoRange, guard.Handler)
	// photoRouter.GET("/first", photoHandler.GetFirstPhoto, guard.Handler)

}
