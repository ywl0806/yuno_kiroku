package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/api/handlers"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
)

type AlbumRouter struct {
	albumHandler handlers.AlbumHandler
}

func NewAlbumRouter(albumHandler handlers.AlbumHandler) *AlbumRouter {
	return &AlbumRouter{albumHandler: albumHandler}
}

func (ar *AlbumRouter) Register(root *echo.Group) {
	album := root.Group("/album")
	guard := middlewares.NewGuard()

	album.GET("/write", ar.albumHandler.GetAlbumsForWrite, guard.Handler)
}
