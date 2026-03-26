package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/handlers"
	"github.com/ywl0806/yuno_kiroku/internal/middlewares"
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
	album.GET("/options", ar.albumHandler.GetAlbumsOptions, guard.Handler)
	album.GET("", ar.albumHandler.GetAllAlbums, guard.Handler)
	album.POST("", ar.albumHandler.CreateAlbum, guard.Handler)
	album.GET("/:id", ar.albumHandler.GetAlbum, guard.Handler)
	album.PUT("/:id", ar.albumHandler.UpdateAlbum, guard.Handler)
	album.DELETE("/:id", ar.albumHandler.DeleteAlbum, guard.Handler)
}
