package middlewares

import (
	"log"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/internal/api/consts"
	"github.com/ywl0806/yuno_kiroku/internal/api/utils/jwt"
)

type AuthUser struct {
	ID          int32  `json:"id" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	GroupId     int32  `json:"group_id" validate:"required"`
	ClanGroupId int32  `json:"clan_group_id" validate:"required"`
}
type Guard struct {
}

func NewGuard() *Guard {
	return &Guard{}
}
func (g *Guard) Handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")

		if token == "" {
			log.Println("token is required")
			return echo.NewHTTPError(401, "Unauthorized")
		}

		claims := jwt.AccessTokenClaims{}

		token, ok := strings.CutPrefix(token, "Bearer ")
		if !ok {
			log.Println("token is required")
			return echo.NewHTTPError(401, "Unauthorized")
		}

		err := jwt.ParseJWT(token, viper.GetString("AUTH_SECRET_KEY"), &claims)
		if err != nil {
			log.Println("token is invalid: ", err)
			return echo.NewHTTPError(401, "Unauthorized")
		}

		authUser := AuthUser{
			ID:          cast.ToInt32(claims.ID),
			Email:       claims.Email,
			GroupId:     cast.ToInt32(claims.GroupId),
			ClanGroupId: cast.ToInt32(claims.ClanGroupId),
		}

		c.Set(consts.AuthUserKey, &authUser)

		return next(c)
	}

}

func GetAuthUser(c echo.Context) *AuthUser {
	return c.Get(consts.AuthUserKey).(*AuthUser)
}
