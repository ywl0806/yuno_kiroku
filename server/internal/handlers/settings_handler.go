package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/handlers/models"
	"github.com/ywl0806/yuno_kiroku/internal/middlewares"
	"github.com/ywl0806/yuno_kiroku/internal/services"
)

type SettingsHandler struct {
	groupService    *services.GroupService
	userService     *services.UserService
	albumService    *services.AlbumService
	kidService      *services.KidService
	identityService *services.IdentityService
}

func NewSettingsHandler(groupService *services.GroupService, userService *services.UserService, albumService *services.AlbumService, kidService *services.KidService, identityService *services.IdentityService) *SettingsHandler {
	return &SettingsHandler{
		groupService:    groupService,
		userService:     userService,
		albumService:    albumService,
		kidService:      kidService,
		identityService: identityService,
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

// @Tags Settings
// @Description Get identity face options
// @Router /settings/identity-face-options [get]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Success 200 {object} []models.IdentityFaceOptionResponse
func (h *SettingsHandler) GetIdentityFaceOptions(c echo.Context) error {
	ctx := c.Request().Context()

	authUser := middlewares.GetAuthUser(c)

	identityFaceOptionRequest := new(models.IdentityFaceOptionRequest)
	if err := c.Bind(identityFaceOptionRequest); err != nil {
		return err
	}
	if err := c.Validate(identityFaceOptionRequest); err != nil {
		return err
	}

	identityFaceOptions, err := h.identityService.FindNewestIdentityFaceImgByFamilyId(ctx, authUser.FamilyId, identityFaceOptionRequest.OnlyNotLinked, identityFaceOptionRequest.WithKidIds, identityFaceOptionRequest.WithUserIds)
	if err != nil {
		return err
	}

	return c.JSON(200, models.NewIdentityFaceOptionResponses(identityFaceOptions))
}
