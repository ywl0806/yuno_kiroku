package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/api/handlers/models"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
	"github.com/ywl0806/yuno_kiroku/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// @Tags User
// @Description create user
// @Router /user [post]
// @Param user body models.CreateUserRequest true "Create User Request"
// @Success 200 {object} map[string]any
func (con *UserHandler) CreateUser(c echo.Context) error {

	var params db.CreateUserParams
	err := models.CreateUserRequest{}.Bind(c, &params)
	if err != nil {
		return err
	}
	ctx := c.Request().Context()

	err = con.userService.ValidateCreateUserParams(ctx, params)
	if err != nil {
		return err
	}
	user, err := con.userService.CreateUser(ctx, params)
	if err != nil {
		return err
	}

	return c.JSON(200, user)

}

// @Tags User
// @Description Get current authenticated user
// @Router /user/me [get]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Success 200 {object} models.MeResponse
func (con *UserHandler) GetMe(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)
	user, err := con.userService.GetMe(c.Request().Context(), authUser.ID)
	if err != nil {
		return err
	}
	return c.JSON(200, models.NewMeResponse(&user))
}

// @Tags User
// @Description Update current user's name
// @Router /user/me [put]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Param body body models.UpdateMeRequest true "Update Me Request"
// @Success 200 {object} models.MeResponse
func (con *UserHandler) UpdateMe(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)
	req := new(models.UpdateMeRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	user, err := con.userService.UpdateMe(c.Request().Context(), authUser.ID, req.Name)
	if err != nil {
		return err
	}
	return c.JSON(200, models.NewMeResponse(&user))
}

// @Tags User
// @Description Update a member's group
// @Router /user/{id}/group [put]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Param id path int true "Member User ID"
// @Param body body models.UpdateMemberGroupRequest true "Update Member Group Request"
// @Success 200 {object} models.MemberResponse
func (con *UserHandler) UpdateMemberGroup(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)
	memberID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	req := new(models.UpdateMemberGroupRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	user, err := con.userService.UpdateMemberGroup(c.Request().Context(), int32(memberID), req.GroupID, authUser.FamilyId)
	if err != nil {
		return err
	}
	return c.JSON(200, models.NewMemberResponse(&db.FindMembersByFamilyIDRow{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		FamilyID: user.FamilyID,
		GroupID:  user.GroupID,
	}))
}

// @Tags User
// @Description Get members in the same family
// @Router /user/members [get]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Success 200 {object} []models.MemberResponse
func (con *UserHandler) GetMembers(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)

	members, err := con.userService.GetMembers(c.Request().Context(), authUser.FamilyId)
	if err != nil {
		return err
	}

	responses := make([]models.MemberResponse, len(members))
	for i, m := range members {
		responses[i] = *models.NewMemberResponse(&m)
	}

	return c.JSON(200, responses)
}
