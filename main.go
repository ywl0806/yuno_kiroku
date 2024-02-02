package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ywl0806/yuno_kiroku/api/route"

	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/ywl0806/yuno_kiroku/api/docs"
)

//go:embed dist
var embededFiles embed.FS

// static file hosting
func getFileSystem() http.FileSystem {

	fsys, err := fs.Sub(embededFiles, "dist")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}
func main() {
	e := echo.New()

	assetHandler := http.FileServer(getFileSystem())
	// swagger setting
	e.GET("/docs/*", echoSwagger.WrapHandler)

	// router setting
	route.Init(e)

	// static file handler
	e.GET("/*", echo.WrapHandler(assetHandler))

	// logger
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status} latency=${latency_human}  ${error}\n ",
	}))
	e.Use(middleware.Recover())

	e.Logger.Fatal(e.Start("localhost:1323"))
}
