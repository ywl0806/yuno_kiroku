package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/handlers"
)

type AuthRouter struct {
	authHandler handlers.AuthHandler
}

func NewAuthRouter(authHandler handlers.AuthHandler) *AuthRouter {
	return &AuthRouter{authHandler: authHandler}
}

func (ar *AuthRouter) Register(root *echo.Group) {
	auth := root.Group("/auth")

	auth.POST("/login", ar.authHandler.Login)

}
