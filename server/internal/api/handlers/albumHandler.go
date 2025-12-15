package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/api/handlers/models"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
	"github.com/ywl0806/yuno_kiroku/internal/api/services"
)

type AlbumHandler struct {
	albumService *services.AlbumService
}

func NewAlbumHandler(albumService *services.AlbumService) *AlbumHandler {
	return &AlbumHandler{albumService: albumService}
}

// @Tags Album
// @Description Get albums for write
// @Router /album/write [get]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Success 200 {object} []models.AlbumResponse
func (con *AlbumHandler) GetAlbumsForWrite(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)
	groupId := authUser.GroupId
	clanGroupId := authUser.ClanGroupId

	albums, err := con.albumService.GetAlbumsForWrite(c.Request().Context(), groupId, clanGroupId)
	if err != nil {
		return err
	}

	albumResponses := make([]models.AlbumResponse, len(albums))
	for i, album := range albums {
		albumResponses[i] = *models.NewAlbumResponse(&album)
	}

	return c.JSON(200, albumResponses)
}
