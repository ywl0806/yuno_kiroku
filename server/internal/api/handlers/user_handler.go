package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/api/handlers/models"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
	"github.com/ywl0806/yuno_kiroku/internal/db"
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
	ctx := c.Request().Context()
	user, err := con.userService.GetUserByID(ctx, authUser.ID)
	if err != nil {
		return err
	}
	isAdmin, err := con.userService.GetGroupIsAdmin(ctx, user.GroupID)
	if err != nil {
		return err
	}
	return c.JSON(200, models.NewMeResponse(&user, isAdmin))
}

// @Tags User
// @Description Update current user's name
// @Router /user/me [put]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Param body body models.UpdateMeRequest true "Update Me Request"
// @Success 200 {object} models.MeResponse
func (con *UserHandler) UpdateMe(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)
	ctx := c.Request().Context()
	req := new(models.UpdateMeRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	user, err := con.userService.UpdateMe(ctx, authUser.ID, req.Name)
	if err != nil {
		return err
	}
	isAdmin, err := con.userService.GetGroupIsAdmin(ctx, user.GroupID)
	if err != nil {
		return err
	}
	return c.JSON(200, models.NewMeResponse(&user, isAdmin))
}

// @Tags User
// @Description Update a member
// @Router /user/{id} [put]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Param id path string true "Member User ID"
// @Param body body models.UpdateMemberRequest true "Update Member Request"
// @Success 200 {object} models.MemberResponse
func (con *UserHandler) UpdateMember(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)
	memberID := c.Param("id")
	if memberID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	req := new(models.UpdateMemberRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	user, err := con.userService.UpdateMember(c.Request().Context(), memberID, authUser.FamilyId, req.GroupID, req.FamilyTitle, req.CustomFamilyTitle)
	if err != nil {
		return err
	}
	return c.JSON(200, models.NewMemberResponse(&user))
}

// @Tags User
// @Description Get members in the same family
// @Router /user/members [get]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Success 200 {object} []models.MemberResponse
func (con *UserHandler) GetMembers(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)

	users, err := con.userService.GetMembers(c.Request().Context(), authUser.FamilyId)
	if err != nil {
		return err
	}

	responses := make([]models.MemberResponse, len(users))
	for i, u := range users {
		responses[i] = *models.NewMemberResponse(&u)
	}

	return c.JSON(200, responses)
}

// @Tags User
// @Description Get a member
// @Router /user/{id} [get]
// @Param Authorization header string true "Authorization" format(bearer) example(bearer token)
// @Param id path string true "Member User ID"
// @Success 200 {object} models.MemberResponse
func (con *UserHandler) GetMember(c echo.Context) error {
	authUser := middlewares.GetAuthUser(c)
	memberID := c.Param("id")
	if memberID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	user, err := con.userService.GetUserByIDAndFamilyID(c.Request().Context(), memberID, authUser.FamilyId)
	if err != nil {
		return err
	}
	return c.JSON(200, models.NewMemberResponse(&user))

}
