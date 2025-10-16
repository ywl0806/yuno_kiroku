package auth

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/internal/api/services/user/store"
	"github.com/ywl0806/yuno_kiroku/internal/api/utils"
)

type AuthController struct {
	userStore *store.UserStore
}

func NewAuthController(
	userStore *store.UserStore,
) *AuthController {
	return &AuthController{
		userStore: userStore,
	}
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
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (con *AuthController) Login(c echo.Context) error {
	var loginRequest struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.Bind(&loginRequest); err != nil {
		fmt.Println("error binding request: ", err)
		return c.JSON(400, map[string]string{"error": "Invalid request"})
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
