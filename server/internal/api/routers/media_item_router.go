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
	mediaItemRouter.GET("/search", mr.mediaItemHandler.SearchMediaItems, guard.Handler)
	mediaItemRouter.GET("/range", mr.mediaItemHandler.GetMediaItemRange, guard.Handler)
	mediaItemRouter.GET("/upload-batch", mr.mediaItemHandler.GetUploadBatches, guard.Handler)
	mediaItemRouter.POST("/upload-batch", mr.mediaItemHandler.CreateUploadBatch, guard.Handler)
	mediaItemRouter.GET("/upload-batch/status", mr.mediaItemHandler.GetUploadBatchStatus, guard.Handler)
	mediaItemRouter.GET("/upload-batch/:id/items", mr.mediaItemHandler.GetUploadBatchItems, guard.Handler)
	mediaItemRouter.POST("/presigned-url", mr.mediaItemHandler.CreatePresignedUpload, guard.Handler)
}
