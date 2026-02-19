package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/handlers/models"
	"github.com/ywl0806/yuno_kiroku/internal/middlewares"
	"github.com/ywl0806/yuno_kiroku/internal/services"
)

type IdentityHandler struct {
	identityService *services.IdentityService
}

func NewIdentityHandler(identityService *services.IdentityService) *IdentityHandler {
	return &IdentityHandler{identityService: identityService}
}

// @Tags Identity
// @Description identity를 수정
// @Router /identity [put]
// @Param updateIdentityRequest body models.UpdateIdentityByIdAndGroupIdRequest true "Update Identity Request"
// @Success 200 {object} models.IdentityResponse
func (h *IdentityHandler) UpdateIdentity(c echo.Context) error {
	updateIdentityRequest := new(models.UpdateIdentityByIdAndGroupIdRequest)
	if err := c.Bind(updateIdentityRequest); err != nil {
		return err
	}
	if err := c.Validate(updateIdentityRequest); err != nil {
		return err
	}

	groupId := middlewares.GetAuthUser(c).GroupId

	err := h.identityService.ValidateUpdateIdentityByIdAndGroupId(c.Request().Context(), updateIdentityRequest.ID, groupId, updateIdentityRequest.Name)
	if err != nil {
		return err
	}
	identity, err := h.identityService.UpdateIdentityByIdAndGroupId(c.Request().Context(), updateIdentityRequest.ID, groupId, updateIdentityRequest.Name)
	if err != nil {
		return err
	}

	return c.JSON(200, models.NewIdentityResponse(&identity))
}

// @Tags Identity
// @Description identity를 조회
// @Router /identity/{id} [get]
// @Param id path string true "Identity ID"
// @Success 200 {object} models.IdentityResponse
func (h *IdentityHandler) FindIdentityByIdAndGroupId(c echo.Context) error {
	findIdentityByIdAndGroupIdRequest := new(models.FindIdentityByIdAndGroupIdRequest)
	if err := c.Bind(findIdentityByIdAndGroupIdRequest); err != nil {
		return err
	}
	if err := c.Validate(findIdentityByIdAndGroupIdRequest); err != nil {
		return err
	}
	groupId := middlewares.GetAuthUser(c).GroupId
	identity, err := h.identityService.FindIdentityByIdAndGroupId(c.Request().Context(), findIdentityByIdAndGroupIdRequest.ID, groupId)
	if err != nil {
		return err
	}
	return c.JSON(200, models.NewIdentityResponse(&identity))
}

// @Tags Identity
// @Description identity 목록을 조회
// @Router /identities [get]
// @Success 200 {object} models.IdentityListResponse
func (h *IdentityHandler) FindIdentities(c echo.Context) error {
	groupId := middlewares.GetAuthUser(c).GroupId
	identities, err := h.identityService.FindIdentitiesByGroupId(c.Request().Context(), groupId)
	if err != nil {
		return err
	}
	return c.JSON(200, models.NewIdentityListResponse(identities))
}
