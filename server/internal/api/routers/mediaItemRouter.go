package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/api/handlers"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
)

type MediaItemRouter struct {
	mediaItemHandler handlers.MediaItemHandler
}

func NewMediaItemRouter(mediaItemHandler handlers.MediaItemHandler) *MediaItemRouter {
	return &MediaItemRouter{mediaItemHandler: mediaItemHandler}
}
func (mr *MediaItemRouter) Register(root *echo.Group) {
	mediaItemRouter := root.Group("/media-item")
	guard := middlewares.NewGuard()

	mediaItemRouter.GET("", mr.mediaItemHandler.GetMediaItems, guard.Handler)
	mediaItemRouter.GET("/range", mr.mediaItemHandler.GetMediaItemRange, guard.Handler)
	mediaItemRouter.POST("/upload", mr.mediaItemHandler.UploadImage, guard.Handler)
}
