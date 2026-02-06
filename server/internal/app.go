package api

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"

	"github.com/labstack/echo/v4/middleware"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/handlers"
	"github.com/ywl0806/yuno_kiroku/internal/middlewares"
	"github.com/ywl0806/yuno_kiroku/internal/providers"
	"github.com/ywl0806/yuno_kiroku/internal/routers"
	"github.com/ywl0806/yuno_kiroku/internal/services"
	"github.com/ywl0806/yuno_kiroku/internal/validator"
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

	storageProvider := providers.NewStorageProvider()
	// Storage service
	storageService := storageProvider.StorageService()

	// service
	userService := services.NewUserService(queries)
	faceService := services.NewFaceService(queries)
	mediaItemService := services.NewMediaItemService(queries, storageService, faceService)
	identityService := services.NewIdentityService(queries)
	albumService := services.NewAlbumService(queries)

	// handler
	userHandler := handlers.NewUserHandler(userService)
	mediaItemHandler := handlers.NewMediaItemHandler(mediaItemService, faceService)
	authHandler := handlers.NewAuthHandler(userService)
	identityHandler := handlers.NewIdentityHandler(identityService)
	albumHandler := handlers.NewAlbumHandler(albumService)

	// root router
	root := e.Group("/api")

	// health check
	root.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	userRouter := routers.NewUserRouter(*userHandler)
	mediaItemRouter := routers.NewMediaItemRouter(*mediaItemHandler)
	authRouter := routers.NewAuthRouter(*authHandler)
	identityRouter := routers.NewIdentityRouter(*identityHandler)
	albumRouter := routers.NewAlbumRouter(*albumHandler)

	e.Validator = validator.NewCustomValidator()

	// 미들웨어 등록 (등록 순서의 역순으로 실행됨)
	// 1. Recover - 패닉 복구 (가장 안쪽에서 실행)
	e.Use(middleware.Recover())
	// 2. Logger - 로깅
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} uri=${uri},\n status=${status},\n latency=${latency_human}\n  ${error}\n ",
	}))
	// 3. ErrorHandler - 에러 처리 (가장 바깥쪽에서 실행, 모든 에러를 처리)
	errorHandler := middlewares.NewErrorHandler()
	e.Use(errorHandler.Handler)

	// init routers
	userRouter.Register(root)
	mediaItemRouter.Register(root)
	authRouter.Register(root)
	identityRouter.Register(root)
	albumRouter.Register(root)

}
