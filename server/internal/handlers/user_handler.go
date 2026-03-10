package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/handlers/models"
	"github.com/ywl0806/yuno_kiroku/internal/middlewares"
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
