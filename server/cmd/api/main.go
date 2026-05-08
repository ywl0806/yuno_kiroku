package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/ywl0806/yuno_kiroku/docs"
	"github.com/ywl0806/yuno_kiroku/internal/api"
	"github.com/ywl0806/yuno_kiroku/pkg/setting"
)

// @BasePath /api
func main() {
	// config
	setting.SettingEnv()
	env := viper.GetString("APP_ENV")

	fmt.Println("env: ", env)
	e := echo.New()

	// /api/* 경로는 API 라우팅
	api.Init(e)

	// swagger docs
	if env == "dev" || env == "local" {
		e.GET("/api/swagger/*", echoSwagger.WrapHandler)
		fmt.Println("Swagger docs: http://localhost:1323/api/swagger/index.html")
	}

	e.Logger.Fatal(e.Start(":1323"))

}
