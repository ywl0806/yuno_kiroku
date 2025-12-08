package handlers

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/internal/api/services"
	"github.com/ywl0806/yuno_kiroku/internal/api/utils"
	"github.com/ywl0806/yuno_kiroku/internal/api/utils/jwt"

	_ "github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	userService *services.UserService
}

func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

type LoginRequest struct {
	Username string `json:"username" validate:"required" example:"admin"`
	Password string `json:"password" validate:"required" example:"password"`
}

// @Description login
//
// @Summary User login
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]string
// @Router /auth/login [post]
func (con *AuthHandler) Login(c echo.Context) error {
	var loginRequest LoginRequest

	if err := c.Bind(&loginRequest); err != nil {
		log.Println("Bind error: ", err)
		return err

	}

	if err := c.Validate(loginRequest); err != nil {
		log.Println("Validate error: ", err)
		return err
	}

	user, _ := con.userService.FindUserByUsername(c.Request().Context(), loginRequest.Username)

	if user.ID == 0 || !utils.CheckPassword(loginRequest.Password, user.Password) {
		return c.JSON(401, map[string]string{"error": "Invalid username or password"})
	}

	accessTokenClaims := &jwt.AccessTokenClaims{
		ID:          cast.ToString(user.ID),
		Email:       user.Username,
		GroupId:     cast.ToString(user.GroupID),
		ClanGroupId: cast.ToString(user.ClanGroupID),
	}
	token, err := jwt.GenerateJWT(accessTokenClaims, viper.GetString("AUTH_SECRET_KEY"), 60*24*30)

	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to generate token"})
	}

	claims := &jwt.AccessTokenClaims{}
	err = jwt.ParseJWT(token, viper.GetString("AUTH_SECRET_KEY"), claims)

	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to generate token"})
	}

	return c.JSON(200,
		map[string]any{
			"message": "Login successful",
			"user": map[string]any{
				"id":            user.ID,
				"username":      user.Username,
				"group_id":      user.GroupID,
				"clan_group_id": user.ClanGroupID,
			},
			"token": token,
		})
}
