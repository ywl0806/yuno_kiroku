package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/handlers"
	"github.com/ywl0806/yuno_kiroku/internal/middlewares"
)

type KidRouter struct {
	kidHandler handlers.KidHandler
}

func NewKidRouter(kidHandler handlers.KidHandler) *KidRouter {
	return &KidRouter{kidHandler: kidHandler}
}

func (r *KidRouter) Register(root *echo.Group) {
	kid := root.Group("/kid")
	guard := middlewares.NewGuard()

	kid.POST("", r.kidHandler.CreateKid, guard.Handler)
	kid.PUT("/:kidId", r.kidHandler.UpdateKid, guard.Handler)
	kid.DELETE("/:kidId", r.kidHandler.DeleteKid, guard.Handler)
}
