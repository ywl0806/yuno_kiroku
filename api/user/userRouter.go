package user

import (
	"github.com/labstack/echo/v4"
)

func Register(root *echo.Group, userController UserController) {
	userRouter := root.Group("/user")

	userRouter.GET("", userController.GetUser)

}
