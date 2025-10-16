package user

import (
	"github.com/labstack/echo/v4"
	groupStore "github.com/ywl0806/yuno_kiroku/internal/api/services/group/store"
	"github.com/ywl0806/yuno_kiroku/internal/api/services/user/models"
	"github.com/ywl0806/yuno_kiroku/internal/api/services/user/store"
	"github.com/ywl0806/yuno_kiroku/internal/api/utils"
)

type UserController struct {
	userStore  *store.UserStore
	groupStore *groupStore.GroupStore
}

func NewUserController(userStore *store.UserStore, groupStore *groupStore.GroupStore) *UserController {
	return &UserController{userStore: userStore, groupStore: groupStore}
}

// @Description get users
// @Router /user [get]
func (con *UserController) GetUsers(c echo.Context) error {

	users, err := con.userStore.FindUsers(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(200, map[string]any{
		"message": "success",
		"data":    users,
	})
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// @Description create user
// @Router /user [post]
// @Param user body CreateUserRequest true "Create User Request"
// @Success 200 {object} map[string]any
func (con *UserController) CreateUser(c echo.Context) error {

	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(400, "Invalid request format")
	}
	name := req.Name
	email := req.Email
	password := req.Password
	if name == "" || email == "" || password == "" {
		return echo.NewHTTPError(400, "Name, email, and password are required")
	}
	hashedPassword, err := utils.HashPassword(password)

	if err != nil {
		return err
	}

	user := models.User{
		Name:     name,
		Email:    &email,
		Password: &hashedPassword,
	}

	err = c.Validate(user)
	if err != nil {
		return err
	}
	group, err := con.groupStore.CreateGroup(c.Request().Context(), "")

	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	user.GroupID = group.ID
	user.IsActive = true

	newUser, err := con.userStore.CreateUser(c.Request().Context(), user)
	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	return c.JSON(200, map[string]any{
		"message": "success",
		"user": map[string]any{
			"name":  newUser.Name,
			"email": newUser.Email,
			"id":    newUser.ID,
		},
	})
}
