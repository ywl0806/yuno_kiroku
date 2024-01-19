package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Hoge struct {
	message string
}

func main() {
	app := echo.New()

	// app.Static("/", "../dist")

	app.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Skipper: nil,
		// Root directory from where the static content is served.
		Root:   "../dist",
		HTML5:  true,
		Browse: false,
	}))

	app.GET("/api", func(ctx echo.Context) error {

		print("hoge")

		return ctx.String(http.StatusOK, "Hello, World!")
	})
	app.Logger.Fatal(app.Start("localhost:1323"))
}
