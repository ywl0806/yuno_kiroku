package main

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"

	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/ywl0806/yuno_kiroku/docs"
	"github.com/ywl0806/yuno_kiroku/internal/api"
	"github.com/ywl0806/yuno_kiroku/internal/api/setting"
)

// @BasePath /api
func main() {
	// config
	setting.SettingEnv()
	mode := viper.GetString("APP_MODE")

	fmt.Println("mode: ", mode)
	e := echo.New()

	// 정적 파일 (업로드된 파일)
	e.Static("/uploads", "uploads")

	// 정적 파일 (React 등) - /api/경로는 제외
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Skipper: func(c echo.Context) bool {
			// /api/경로는 정적 파일 제공 제외
			return strings.HasPrefix(c.Path(), "/api/")
		},
		HTML5: true,
		Root:  "dist", // React 빌드 결과물이 들어 있는 디렉토리
	}))

	// /swagger/* 경로는 swagger docs
	if mode == "dev" {
		e.GET("/api/swagger/*", echoSwagger.WrapHandler)
	}

	// /api/* 경로는 API 라우팅
	api.Init(e)
	e.GET("/*", func(c echo.Context) error {
		path := c.Request().URL.Path
		if strings.HasPrefix(path, "/api/") {
			return echo.NewHTTPError(404, "Not Found")
		}
		return c.File("dist/index.html")
	})
	// 로거 & 에러 복구
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} uri=${uri},\n status=${status},\n latency=${latency_human}\n  ${error}\n ",
	}))
	e.Use(middleware.Recover())

	e.Logger.Fatal(e.Start(":1323"))
}
