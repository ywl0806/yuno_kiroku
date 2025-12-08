package handlers

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/api/utils"
	"github.com/ywl0806/yuno_kiroku/internal/db"

	"github.com/ywl0806/yuno_kiroku/internal/api/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// // @Description get users
// // @Router /user [get]
// func (con *UserHandler) GetUsers(c echo.Context) error {

// 	users, err := con.userStore.FindUsers(c.Request().Context())
// 	if err != nil {
// 		return err
// 	}

// 	return c.JSON(200, map[string]any{
// 		"message": "success",
// 		"data":    users,
// 	})
// }

type CreateUserRequest struct {
	Name        string `json:"name" validate:"required"`
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required,min=6"`
	GroupID     int32  `json:"group_id" validate:"required"`
	ClanGroupID int32  `json:"clan_group_id" validate:"required"`
}

func (CreateUserRequest) bind(c echo.Context, params *db.CreateUserParams) error {
	req := new(CreateUserRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}
	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}
	params.Name = sql.NullString{String: req.Name, Valid: true}
	params.Username = req.Username
	params.Password = hashPassword
	params.GroupID = req.GroupID
	params.ClanGroupID = req.ClanGroupID
	return nil
}

// @Tags User
// @Description create user
// @Router /user [post]
// @Param user body CreateUserRequest true "Create User Request"
// @Success 200 {object} map[string]any
func (con *UserHandler) CreateUser(c echo.Context) error {

	var params db.CreateUserParams
	err := CreateUserRequest{}.bind(c, &params)
	if err != nil {
		return err
	}
	ctx := c.Request().Context()

	err = con.userService.ValidateCreateUserParams(ctx, params)
	if err != nil {
		return HandleServiceError(err)
	}
	user, err := con.userService.CreateUser(ctx, params)
	if err != nil {
		return HandleServiceError(err)
	}

	return c.JSON(200, user)

}
