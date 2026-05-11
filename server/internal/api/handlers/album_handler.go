package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/api/handlers/models"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
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

// @Tags Album
// @Description Get all albums for the family (settings use)
// @Router /album [get]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Success 200 {object} []models.AlbumResponse
func (con *AlbumHandler) GetAllAlbums(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)
	albums, err := con.albumService.GetAllAlbums(c.Request().Context(), authUser.FamilyId)
	if err != nil {
		return err
	}
	responses := make([]models.AlbumResponse, len(albums))
	for i, a := range albums {
		responses[i] = *models.NewAlbumResponse(&a)
	}
	return c.JSON(200, responses)
}

// @Tags Album
// @Description Create a new album
// @Router /album [post]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Param body body models.CreateAlbumRequest true "Create Album Request"
// @Success 201 {object} models.AlbumWithPermissionsResponse
func (con *AlbumHandler) CreateAlbum(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)
	req := new(models.CreateAlbumRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}
	perms := make([]services.GroupPermission, len(req.Permissions))
	for i, p := range req.Permissions {
		perms[i] = services.GroupPermission{GroupID: p.GroupID, Permission: p.Permission}
	}
	album, err := con.albumService.CreateAlbum(c.Request().Context(), authUser.FamilyId, req.Name, req.IsCommon, perms)
	if err != nil {
		return err
	}
	albumPerms, err := con.albumService.GetAlbumPermissions(c.Request().Context(), album.ID, authUser.FamilyId)
	if err != nil {
		return err
	}
	return c.JSON(201, models.NewAlbumWithPermissionsResponse(&album, albumPerms))
}

// @Tags Album
// @Description Get album with permissions
// @Router /album/:id [get]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Param id path string true "Album ID"
// @Success 200 {object} models.AlbumWithPermissionsResponse
func (con *AlbumHandler) GetAlbum(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)
	id := c.Param("id")
	album, perms, err := con.albumService.GetAlbumWithPermissions(c.Request().Context(), id, authUser.FamilyId)
	if err != nil {
		return err
	}
	return c.JSON(200, models.NewAlbumWithPermissionsResponse(&album, perms))
}

// @Tags Album
// @Description Update an album
// @Router /album/:id [put]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Param id path string true "Album ID"
// @Param body body models.UpdateAlbumRequest true "Update Album Request"
// @Success 200 {object} models.AlbumWithPermissionsResponse
func (con *AlbumHandler) UpdateAlbum(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)
	id := c.Param("id")
	req := new(models.UpdateAlbumRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}
	perms := make([]services.GroupPermission, len(req.Permissions))
	for i, p := range req.Permissions {
		perms[i] = services.GroupPermission{GroupID: p.GroupID, Permission: p.Permission}
	}
	album, err := con.albumService.UpdateAlbum(c.Request().Context(), id, authUser.FamilyId, req.Name, perms)
	if err != nil {
		return err
	}
	albumPerms, err := con.albumService.GetAlbumPermissions(c.Request().Context(), album.ID, authUser.FamilyId)
	if err != nil {
		return err
	}
	return c.JSON(200, models.NewAlbumWithPermissionsResponse(&album, albumPerms))
}

// @Tags Album
// @Description Delete an album
// @Router /album/:id [delete]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Param id path string true "Album ID"
// @Success 204
func (con *AlbumHandler) DeleteAlbum(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)
	id := c.Param("id")
	if err := con.albumService.DeleteAlbum(c.Request().Context(), id, authUser.FamilyId); err != nil {
		return err
	}
	return c.NoContent(204)
}
