package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	worker "github.com/ywl0806/yuno_kiroku/internal/worker"
	"github.com/ywl0806/yuno_kiroku/pkg/setting"
)

func main() {
	setting.SettingEnv()
	port := viper.GetString("PORT")
	if port == "" {
		port = "1325"
	}

	e := echo.New()
	worker.InitLocal(e)

	fmt.Printf("Resize Worker 시작: :%s\n", port)
	e.Logger.Fatal(e.Start(":" + port))
}
