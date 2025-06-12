package user

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/api/services/user/models"
	"github.com/ywl0806/yuno_kiroku/internal/api/services/user/store"
	"github.com/ywl0806/yuno_kiroku/internal/api/utils"
)

type UserController struct {
	userStore *store.UserStore
}

func NewUserController(userStore *store.UserStore) *UserController {
	return &UserController{userStore: userStore}
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

// @Description create user
// @Router /user [post]
// @Param name formData string true "name"
// @Param email formData string true "email" format(email)
// @Param password formData string true "password"
func (con *UserController) CreateUser(c echo.Context) error {

	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")
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
