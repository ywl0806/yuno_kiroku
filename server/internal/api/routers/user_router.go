package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/api/handlers"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
)

type UserRouter struct {
	userHandler handlers.UserHandler
}

func NewUserRouter(userHandler handlers.UserHandler) *UserRouter {
	return &UserRouter{userHandler: userHandler}
}

func (ur *UserRouter) Register(root *echo.Group) {
	userRouter := root.Group("/user")
	guard := middlewares.NewGuard()

	userRouter.POST("", ur.userHandler.CreateUser)
	userRouter.GET("/me", ur.userHandler.GetMe, guard.Handler)
	userRouter.PUT("/me", ur.userHandler.UpdateMe, guard.Handler)
	userRouter.GET("/members", ur.userHandler.GetMembers, guard.Handler)
	userRouter.PUT("/:id", ur.userHandler.UpdateMember, guard.Handler)
	userRouter.GET("/:id", ur.userHandler.GetMember, guard.Handler)
}
