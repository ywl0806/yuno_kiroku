package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/api/handlers/models"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
	"github.com/ywl0806/yuno_kiroku/internal/services"
)

type IdentityHandler struct {
	identityService *services.IdentityService
}

func NewIdentityHandler(identityService *services.IdentityService) *IdentityHandler {
	return &IdentityHandler{identityService: identityService}
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
	familyId := middlewares.GetAuthUser(c).FamilyId
	identity, err := h.identityService.FindIdentityByIdAndFamilyId(c.Request().Context(), findIdentityByIdAndGroupIdRequest.ID, familyId)
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
	familyId := middlewares.GetAuthUser(c).FamilyId
	identities, err := h.identityService.FindIdentitiesByFamilyId(c.Request().Context(), familyId)
	if err != nil {
		return err
	}
	return c.JSON(200, models.NewIdentityListResponse(identities))
}

// @Tags Identity
// @Description identity 옵션 조회
// @Router /identity/options [get]
// @Success 200 {object} models.IdentityOptionResponse
func (h *IdentityHandler) GetIdentityOptions(c echo.Context) error {
	familyId := middlewares.GetAuthUser(c).FamilyId
	identityOptions, err := h.identityService.GetIdentityOptions(c.Request().Context(), familyId)
	if err != nil {
		return err
	}
	return c.JSON(200, models.NewIdentityOptionResponses(identityOptions))
}
