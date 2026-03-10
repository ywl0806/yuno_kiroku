package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/handlers/models"
	"github.com/ywl0806/yuno_kiroku/internal/middlewares"
	"github.com/ywl0806/yuno_kiroku/internal/services"
)

type GroupHandler struct {
	groupService *services.GroupService
}

func NewGroupHandler(groupService *services.GroupService) *GroupHandler {
	return &GroupHandler{groupService: groupService}
}

// @Tags Group
// @Description Get album groups for the authenticated user's family
// @Router /group [get]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Success 200 {object} []models.GroupResponse
func (h *GroupHandler) GetGroups(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)

	groups, err := h.groupService.GetGroups(c.Request().Context(), authUser.FamilyId)
	if err != nil {
		return err
	}

	responses := make([]models.GroupResponse, len(groups))
	for i, g := range groups {
		responses[i] = *models.NewGroupResponse(&g)
	}

	return c.JSON(200, responses)
}
