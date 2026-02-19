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

	// 소셜 로그인: 로그인 페이지로 리다이렉트
	auth.GET("/line", ar.authHandler.LineLoginRedirect)
	auth.GET("/line/callback", ar.authHandler.LineCallback)
	auth.GET("/kakao", ar.authHandler.KakaoLoginRedirect)
	auth.GET("/kakao/callback", ar.authHandler.KakaoCallback)
}
