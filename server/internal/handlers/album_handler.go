package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/handlers/models"
	"github.com/ywl0806/yuno_kiroku/internal/middlewares"
	"github.com/ywl0806/yuno_kiroku/internal/services"
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
	familyId := authUser.FamilyId
	groupId := authUser.GroupId

	albums, err := con.albumService.GetAlbumsForWrite(c.Request().Context(), familyId, groupId)
	if err != nil {
		return err
	}

	albumResponses := make([]models.AlbumResponse, len(albums))
	for i, album := range albums {
		albumResponses[i] = *models.NewAlbumResponse(&album)
	}

	return c.JSON(200, albumResponses)
}

// @Tags Album
// @Description Get albums options
// @Router /album/options [get]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Success 200 {object} []models.AlbumOptionResponse
func (con *AlbumHandler) GetAlbumsOptions(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)
	familyId := authUser.FamilyId
	groupId := authUser.GroupId
	opts, err := con.albumService.GetAlbumsOptions(c.Request().Context(), familyId, groupId)
	if err != nil {
		return err
	}
	return c.JSON(200, models.NewAlbumOptionResponses(opts))
}
