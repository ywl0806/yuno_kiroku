package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/handlers"
	"github.com/ywl0806/yuno_kiroku/internal/middlewares"
)

type SettingsRouter struct {
	settingsHandler handlers.SettingsHandler
}

func NewSettingsRouter(settingsHandler handlers.SettingsHandler) *SettingsRouter {
	return &SettingsRouter{settingsHandler: settingsHandler}
}

func (r *SettingsRouter) Register(root *echo.Group) {
	group := root.Group("/settings")
	guard := middlewares.NewGuard()

	group.GET("/data", r.settingsHandler.GetSettingsData, guard.Handler)
}
