package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ywl0806/yuno_kiroku/api/route"
)

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
	route.New(app)

	app.Logger.Fatal(app.Start("localhost:1323"))
}
