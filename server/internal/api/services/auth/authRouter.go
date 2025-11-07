package auth

import "github.com/labstack/echo/v4"

func Register(root *echo.Group, authHandler AuthHandler) {
	auth := root.Group("/auth")

	auth.POST("/login", authHandler.Login)

}
