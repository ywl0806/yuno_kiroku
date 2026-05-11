package middlewares

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/internal/consts"
	"github.com/ywl0806/yuno_kiroku/internal/utils/jwt"
)

type AuthUser struct {
	ID       int32  `json:"id" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	FamilyId int32  `json:"family_id" validate:"required"`
	GroupId  int32  `json:"group_id" validate:"required"`
}

type Guard struct{}

func NewGuard() *Guard {
	return &Guard{}
}

func (g *Guard) Handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(consts.AccessTokenCookieName)
		if err != nil || cookie.Value == "" {
			log.Println("access token cookie is required")
			return echo.NewHTTPError(401, "Unauthorized")
		}

		claims := jwt.AccessTokenClaims{}
		if err := jwt.ParseJWT(cookie.Value, viper.GetString("AUTH_SECRET_KEY"), &claims); err != nil {
			log.Println("token is invalid: ", err)
			return echo.NewHTTPError(401, "Unauthorized")
		}

		c.Set(consts.AuthUserKey, &AuthUser{
			ID:       cast.ToInt32(claims.ID),
			Email:    claims.Email,
			FamilyId: cast.ToInt32(claims.FamilyId),
			GroupId:  cast.ToInt32(claims.GroupId),
		})

		return next(c)
	}
}

func GetAuthUser(c echo.Context) *AuthUser {
	return c.Get(consts.AuthUserKey).(*AuthUser)
}
