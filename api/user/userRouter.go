package user

import (
	"github.com/labstack/echo/v4"
)

func Register(root *echo.Group, userHandler UserContoller) {
	userRouter := root.Group("/user")

	userRouter.GET("", userHandler.GetUser)

}
