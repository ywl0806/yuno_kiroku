package user

import (
	"github.com/labstack/echo/v4"
)

func Route(root *echo.Group) *echo.Group {
	g := root.Group("/user")
	g.GET("/hi", func(c echo.Context) error {

		return c.String(200, "hello")
	})
	return g
}
