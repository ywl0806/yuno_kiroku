package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/api/handlers"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
)

type GroupRouter struct {
	groupHandler handlers.GroupHandler
}

func NewGroupRouter(groupHandler handlers.GroupHandler) *GroupRouter {
	return &GroupRouter{groupHandler: groupHandler}
}

func (r *GroupRouter) Register(root *echo.Group) {
	group := root.Group("/group")
	guard := middlewares.NewGuard()

	group.GET("", r.groupHandler.GetGroups, guard.Handler)
	group.POST("", r.groupHandler.CreateGroup, guard.Handler)
	group.PUT("/:id", r.groupHandler.UpdateGroup, guard.Handler)
	group.DELETE("/:id", r.groupHandler.DeleteGroup, guard.Handler)
}
