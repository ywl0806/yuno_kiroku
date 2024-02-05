package main

import (
	"embed"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ywl0806/yuno_kiroku/api/route"
	"github.com/ywl0806/yuno_kiroku/api/setting"

	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/ywl0806/yuno_kiroku/api/docs"
)

//go:embed dist
var webAssets embed.FS

// @BasePath /api
func main() {
	setting.SettingEnv()

	e := echo.New()

	// static file handler
	e.Static("/uploads", "uploads")

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Skipper: func(c echo.Context) bool {
			basePath := c.Path()[1:8]
			return basePath == "swagger"
		},
		HTML5:      true,
		Root:       "dist",
		Filesystem: http.FS(webAssets),
	}))

	// swagger setting
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// router setting
	route.Init(e)

	// logger
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status} latency=${latency_human}  ${error}\n ",
	}))
	e.Use(middleware.Recover())

	e.Logger.Fatal(e.Start("localhost:1323"))
}
