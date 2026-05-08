package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/api/handlers"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
)

type InviteRouter struct {
	inviteHandler handlers.InviteHandler
}

func NewInviteRouter(inviteHandler handlers.InviteHandler) *InviteRouter {
	return &InviteRouter{inviteHandler: inviteHandler}
}

func (ir *InviteRouter) Register(root *echo.Group) {
	invite := root.Group("/invite")
	guard := middlewares.NewGuard()

	invite.GET("/validate", ir.inviteHandler.ValidateInvite)
	invite.POST("", ir.inviteHandler.CreateInvite, guard.Handler)
}
