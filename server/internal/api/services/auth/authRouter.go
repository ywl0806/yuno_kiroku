package auth

import "github.com/labstack/echo/v4"

func Register(root *echo.Group, authController AuthController) {
	authRouter := root.Group("/auth")

	authRouter.POST("/login", authController.Login)

}
