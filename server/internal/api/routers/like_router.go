package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/api/handlers"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
)

type LikeRouter struct {
	likeHandler    handlers.LikeHandler
	tagHandler     handlers.TagHandler
}

func NewLikeRouter(likeHandler handlers.LikeHandler, tagHandler handlers.TagHandler) *LikeRouter {
	return &LikeRouter{
		likeHandler: likeHandler,
		tagHandler:  tagHandler,
	}
}

func (r *LikeRouter) Register(root *echo.Group) {
	guard := middlewares.NewGuard()

	// 좋아요 관련 라우트 (/media-item/:id/like)
	mediaItemRouter := root.Group("/media-item")
	mediaItemRouter.GET("/:id/like", r.likeHandler.IsMediaItemLiked, guard.Handler)
	mediaItemRouter.POST("/:id/like", r.likeHandler.LikeMediaItem, guard.Handler)
	mediaItemRouter.DELETE("/:id/like", r.likeHandler.UnlikeMediaItem, guard.Handler)

	// 미디어 아이템 태그 관련 라우트 (/media-item/:id/tag)
	mediaItemRouter.GET("/:id/tag", r.tagHandler.GetTagsForMediaItem, guard.Handler)
	mediaItemRouter.POST("/:id/tag", r.tagHandler.AddTagToMediaItem, guard.Handler)
	mediaItemRouter.DELETE("/:id/tag/:tagId", r.tagHandler.RemoveTagFromMediaItem, guard.Handler)

	// 태그 CRUD (/tag)
	tagRouter := root.Group("/tag")
	tagRouter.GET("", r.tagHandler.GetTags, guard.Handler)
	tagRouter.POST("", r.tagHandler.CreateTag, guard.Handler)
	tagRouter.DELETE("/:id", r.tagHandler.DeleteTag, guard.Handler)
}
