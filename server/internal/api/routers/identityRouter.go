package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/api/handlers"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
)

type IdentityRouter struct {
	identityHandler handlers.IdentityHandler
}

func NewIdentityRouter(identityHandler handlers.IdentityHandler) *IdentityRouter {
	return &IdentityRouter{identityHandler: identityHandler}
}

func (ir *IdentityRouter) Register(root *echo.Group) {
	identityRouter := root.Group("/identity")
	guard := middlewares.NewGuard()

	identityRouter.PUT("", ir.identityHandler.UpdateIdentity, guard.Handler)
	identityRouter.GET("/:id", ir.identityHandler.FindIdentityByIdAndGroupId, guard.Handler)
	identityRouter.GET("", ir.identityHandler.FindIdentities, guard.Handler)
}
