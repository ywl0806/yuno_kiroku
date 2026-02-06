package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/handlers"
)

type UserRouter struct {
	userHandler handlers.UserHandler
}

func NewUserRouter(userHandler handlers.UserHandler) *UserRouter {
	return &UserRouter{userHandler: userHandler}
}

func (ur *UserRouter) Register(root *echo.Group) {
	userRouter := root.Group("/user")

	userRouter.POST("", ur.userHandler.CreateUser)

}
