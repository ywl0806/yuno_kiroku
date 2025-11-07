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
	"github.com/ywl0806/yuno_kiroku/pkg/setting"
)

// @BasePath /api
func main() {
	// config
	setting.SettingEnv()
	mode := viper.GetString("APP_MODE")

	fmt.Println("mode: ", mode)
	e := echo.New()

	// /api/* 경로는 API 라우팅
	api.Init(e)

	// /swagger/* 경로는 swagger docs
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

	// swagger docs
	e.GET("/api/swagger/*", echoSwagger.WrapHandler)

	fmt.Println("Swagger docs: http://localhost:1323/api/swagger/index.html")
	e.GET("/*", func(c echo.Context) error {
		path := c.Request().URL.Path
		if strings.HasPrefix(path, "/api/") {
			return echo.NewHTTPError(404, "Not Found")
		}
		return c.File("dist/index.html")
	})
	e.Logger.Fatal(e.Start(":1323"))

}
