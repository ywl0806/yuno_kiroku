package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func main() {
	app := echo.New()

	app.Static("/", "../dist")
	app.GET("/api", func(ctx echo.Context) error {
		print("hoge")
		return ctx.String(http.StatusOK, "Hello, World!")
	})
	app.Logger.Fatal(app.Start("localhost:1323"))
}
