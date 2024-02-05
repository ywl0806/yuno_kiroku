package main

import (
	"embed"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/ywl0806/yuno_kiroku/api"
	"github.com/ywl0806/yuno_kiroku/api/setting"

	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/ywl0806/yuno_kiroku/docs"
)

//go:embed dist
var webAssets embed.FS

// @BasePath /api
func main() {
	// config
	setting.SettingEnv()

	e := echo.New()

	// static file
	e.Static("/uploads", "uploads")
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Path(), "/swagger/")
		},
		HTML5:      true,
		Root:       "dist",
		Filesystem: http.FS(webAssets),
	}))

	// swagger setting
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// api setting
	api.Init(e)

	// logger
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status} latency=${latency_human}  ${error}\n ",
	}))
	e.Use(middleware.Recover())

	e.Logger.Fatal(e.Start("localhost:1323"))
}
