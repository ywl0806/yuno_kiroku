package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/internal/api/services/user/store"
	"github.com/ywl0806/yuno_kiroku/internal/api/utils"
)

type AuthHandler struct {
	userStore *store.UserStore
}

func NewAuthHandler(userStore *store.UserStore) *AuthHandler {
	return &AuthHandler{userStore: userStore}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
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

	if err := c.Validate(loginRequest); err != nil {
		return err
	}

	user, err := con.userStore.FindUserByEmail(c.Request().Context(), loginRequest.Email)

	if err != nil {
		return err
	}

	if user.Password == nil || !utils.CheckPassword(loginRequest.Password, *user.Password) {
		return c.JSON(401, map[string]string{"error": "Invalid credentials"})
	}

	accessTokenClaims := &utils.AccessTokenClaims{
		ID:      user.ID.Hex(),
		Email:   *user.Email,
		GroupId: user.GroupID.Hex(),
	}
	if user.ClanGroupId != nil {
		accessTokenClaims.ClanGroupId = user.ClanGroupId.Hex()
	} else {
		accessTokenClaims.ClanGroupId = ""
	}
	token, err := utils.GenerateJWT(accessTokenClaims, viper.GetString("AUTH_SECRET_KEY"), 60*24*30)

	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to generate token"})
	}

	claims := &utils.AccessTokenClaims{}
	err = utils.ParseJWT(token, viper.GetString("AUTH_SECRET_KEY"), claims)

	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to generate token"})
	}

	return c.JSON(200,
		map[string]string{
			"message": "Login successful",
			"email":   loginRequest.Email,
			"token":   token,
		})
}
