package handlers

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/api/handlers/models"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
	"github.com/ywl0806/yuno_kiroku/internal/services"
)

type GroupHandler struct {
	groupService services.GroupService
}

func NewGroupHandler(groupService services.GroupService) *GroupHandler {
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

// @Tags Group
// @Description Create a new album group
// @Router /group [post]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Param body body models.CreateGroupRequest true "Create Group Request"
// @Success 201 {object} models.GroupResponse
func (h *GroupHandler) CreateGroup(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)
	req := new(models.CreateGroupRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}
	group, err := h.groupService.CreateGroup(c.Request().Context(), authUser.FamilyId, req.Name)
	if err != nil {
		return err
	}
	return c.JSON(201, models.NewGroupResponse(&group))
}

// @Tags Group
// @Description Update an album group
// @Router /group/:id [put]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Param id path int true "Group ID"
// @Param body body models.UpdateGroupRequest true "Update Group Request"
// @Success 200 {object} models.GroupResponse
func (h *GroupHandler) UpdateGroup(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(400, "invalid group id")
	}
	req := new(models.UpdateGroupRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}
	group, err := h.groupService.UpdateGroup(c.Request().Context(), int32(id), authUser.FamilyId, req.Name)
	if err != nil {
		return err
	}
	return c.JSON(200, models.NewGroupResponse(&group))
}

// @Tags Group
// @Description Delete an album group
// @Router /group/:id [delete]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Param id path int true "Group ID"
// @Success 204
func (h *GroupHandler) DeleteGroup(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(400, "invalid group id")
	}
	if err := h.groupService.DeleteGroup(c.Request().Context(), int32(id), authUser.FamilyId); err != nil {
		return err
	}
	return c.NoContent(204)
}
