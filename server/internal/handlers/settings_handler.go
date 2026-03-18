package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/handlers/models"
	"github.com/ywl0806/yuno_kiroku/internal/middlewares"
	"github.com/ywl0806/yuno_kiroku/internal/services"
)

type SettingsHandler struct {
	groupService *services.GroupService
	userService  *services.UserService
	albumService *services.AlbumService
	kidService   *services.KidService
}

func NewSettingsHandler(groupService *services.GroupService, userService *services.UserService, albumService *services.AlbumService, kidService *services.KidService) *SettingsHandler {
	return &SettingsHandler{
		groupService: groupService,
		userService:  userService,
		albumService: albumService,
		kidService:   kidService,
	}
}

// @Tags Settings
// @Description Get all settings data (groups, members, albums) in a single request
// @Router /settings/data [get]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Success 200 {object} models.SettingsDataResponse
func (h *SettingsHandler) GetSettingsData(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)
	ctx := c.Request().Context()

	groups, err := h.groupService.GetGroups(ctx, authUser.FamilyId)
	if err != nil {
		return err
	}

	members, err := h.userService.GetMembers(ctx, authUser.FamilyId)
	if err != nil {
		return err
	}

	albums, err := h.albumService.GetAlbumsForWrite(ctx, authUser.FamilyId, authUser.GroupId)
	if err != nil {
		return err
	}

	kids, err := h.kidService.GetKids(ctx, authUser.FamilyId)
	if err != nil {
		return err
	}

	groupResponses := make([]models.GroupResponse, len(groups))
	for i, g := range groups {
		groupResponses[i] = *models.NewGroupResponse(&g)
	}

	memberResponses := make([]models.MemberResponse, len(members))
	for i, m := range members {
		memberResponses[i] = *models.NewMemberResponse(&m)
	}

	albumResponses := make([]models.AlbumResponse, len(albums))
	for i, a := range albums {
		albumResponses[i] = *models.NewAlbumResponse(&a)
	}

	kidResponses := make([]models.KidResponse, len(kids))
	for i, k := range kids {
		kidResponses[i] = *models.NewKidResponse(&k)
	}

	return c.JSON(200, models.SettingsDataResponse{
		Groups:  groupResponses,
		Members: memberResponses,
		Albums:  albumResponses,
		Kids:    kidResponses,
	})
}
