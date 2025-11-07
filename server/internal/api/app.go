package api

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/pkg/storage"

	_ "github.com/lib/pq"

	"github.com/labstack/echo/v4/middleware"
	"github.com/ywl0806/yuno_kiroku/internal/api/services/auth"
	"github.com/ywl0806/yuno_kiroku/internal/api/services/photo"
	"github.com/ywl0806/yuno_kiroku/internal/api/services/user"
	"github.com/ywl0806/yuno_kiroku/internal/api/validator"
	"github.com/ywl0806/yuno_kiroku/internal/db"
)

// Initialize the root router on the app
func Init(e *echo.Echo) {

	dbTx, err := sql.Open("postgres", viper.GetString("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	queries := db.New(dbTx)

	// Storage service
	sStorage := storage.NewLocalStorageService("standard")
	lStorage := storage.NewLocalStorageService("longterm")

	userService := user.NewUserService(queries)
	photoService := photo.NewPhotoService(sStorage, lStorage)

	userHandler := user.NewUserHandler(userService)
	photoHandler := photo.NewPhotoHandler(photoService)
	authHandler := auth.NewAuthHandler(userService)

	// root router
	root := e.Group("/api")

	// health check
	root.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})
	// init routers
	user.Register(root, *userHandler)
	photo.Register(root, *photoHandler)
	auth.Register(root, *authHandler)

	e.Validator = validator.NewCustomValidator()
	// 로거 & 에러 복구
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} uri=${uri},\n status=${status},\n latency=${latency_human}\n  ${error}\n ",
	}))
	e.Use(middleware.Recover())
}
