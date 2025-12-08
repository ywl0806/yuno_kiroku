package api

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/pkg/storage"

	_ "github.com/lib/pq"

	"github.com/labstack/echo/v4/middleware"
	"github.com/ywl0806/yuno_kiroku/internal/api/handlers"
	"github.com/ywl0806/yuno_kiroku/internal/api/routers"
	"github.com/ywl0806/yuno_kiroku/internal/api/services"
	"github.com/ywl0806/yuno_kiroku/internal/api/validator"
	"github.com/ywl0806/yuno_kiroku/internal/db"
)

// Initialize the root router on the app
func Init(e *echo.Echo) {

	dbTx, err := sql.Open("postgres", viper.GetString("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	// 연결 테스트
	if err := dbTx.Ping(); err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	queries := db.New(dbTx)

	// Storage service
	sStorage := storage.NewLocalStorageService("standard")
	lStorage := storage.NewLocalStorageService("longterm")

	userService := services.NewUserService(queries)
	photoService := services.NewPhotoService(queries, sStorage, lStorage)
	faceService := services.NewFaceService(queries)
	identityService := services.NewIdentityService(queries)

	userHandler := handlers.NewUserHandler(userService)
	photoHandler := handlers.NewPhotoHandler(photoService, faceService)
	authHandler := handlers.NewAuthHandler(userService)
	identityHandler := handlers.NewIdentityHandler(identityService)

	// root router
	root := e.Group("/api")

	// health check
	root.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	userRouter := routers.NewUserRouter(*userHandler)
	photoRouter := routers.NewPhotoRouter(*photoHandler)
	authRouter := routers.NewAuthRouter(*authHandler)
	identityRouter := routers.NewIdentityRouter(*identityHandler)

	// init routers
	userRouter.Register(root)
	photoRouter.Register(root)
	authRouter.Register(root)
	identityRouter.Register(root)

	e.Validator = validator.NewCustomValidator()
	// 로거 & 에러 복구
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} uri=${uri},\n status=${status},\n latency=${latency_human}\n  ${error}\n ",
	}))
	e.Use(middleware.Recover())
}
