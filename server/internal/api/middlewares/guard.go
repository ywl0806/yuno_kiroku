package middlewares

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/internal/api/consts"
	"github.com/ywl0806/yuno_kiroku/internal/api/utils"
)

type Guard struct {
}

func NewGuard() *Guard {
	return &Guard{}
}
func (g *Guard) Handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return echo.NewHTTPError(401, "Unauthorized")
		}

		claims := utils.AccessTokenClaims{}

		token, ok := strings.CutPrefix(token, "Bearer ")
		if !ok {
			return echo.NewHTTPError(401, "Unauthorized")
		}

		err := utils.ParseJWT(token, viper.GetString("AUTH_SECRET_KEY"), &claims)
		if err != nil {
			return echo.NewHTTPError(401, "Unauthorized")
		}

		c.Set(consts.UserIdKey, claims.ID)
		c.Set(consts.UserEmailKey, claims.Email)
		c.Set(consts.UserGroupIdKey, claims.GroupId)
		c.Set(consts.UserRoleKey, claims.Role)

		return next(c)
	}

}
