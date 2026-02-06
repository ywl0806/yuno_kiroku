package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/internal/consts"
	"github.com/ywl0806/yuno_kiroku/internal/services"
	"github.com/ywl0806/yuno_kiroku/internal/utils"
	"github.com/ywl0806/yuno_kiroku/internal/utils/jwt"

	_ "github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	userService   *services.UserService
	authSecretKey string
}

func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{userService: userService, authSecretKey: viper.GetString("AUTH_SECRET_KEY")}
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
	token, err := jwt.GenerateJWT(accessTokenClaims, con.authSecretKey, consts.AccessTokenCookieMaxAge)

	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to generate token"})
	}

	refreshTokenClaims := &jwt.RefreshTokenClaims{
		ID: cast.ToString(user.ID),
	}

	refreshToken, err := jwt.GenerateJWT(refreshTokenClaims, con.authSecretKey, consts.RefreshTokenCookieMaxAge)

	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to generate token"})
	}

	c.SetCookie(&http.Cookie{
		Name:     consts.RefreshTokenCookieName,
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		MaxAge:   consts.RefreshTokenCookieMaxAge,
	})

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
