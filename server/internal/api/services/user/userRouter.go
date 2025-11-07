package user

import (
	"github.com/labstack/echo/v4"
)

func Register(root *echo.Group, userHandler UserHandler) {
	userRouter := root.Group("/user")

	userRouter.POST("", userHandler.CreateUser)

}
