package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/api/handlers"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
)

type PhotoRouter struct {
	photoHandler handlers.PhotoHandler
}

func NewPhotoRouter(photoHandler handlers.PhotoHandler) *PhotoRouter {
	return &PhotoRouter{photoHandler: photoHandler}
}
func (pr *PhotoRouter) Register(root *echo.Group) {
	photoRouter := root.Group("/photo")
	guard := middlewares.NewGuard()

	photoRouter.GET("", pr.photoHandler.GetPhotos, guard.Handler)
	photoRouter.POST("/upload", pr.photoHandler.UploadPhoto, guard.Handler)
	photoRouter.POST("/upload-live", pr.photoHandler.UploadLivePhoto, guard.Handler)
	photoRouter.GET("/range", pr.photoHandler.GetPhotoRange, guard.Handler)
	photoRouter.GET("/identity/:identity_id/random", pr.photoHandler.GetIdentityRandomPhoto, guard.Handler)
}
