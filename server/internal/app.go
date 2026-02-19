package api

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"

	"github.com/labstack/echo/v4/middleware"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/handlers"
	"github.com/ywl0806/yuno_kiroku/internal/i18n"
	"github.com/ywl0806/yuno_kiroku/internal/middlewares"
	"github.com/ywl0806/yuno_kiroku/internal/oauth"
	"github.com/ywl0806/yuno_kiroku/internal/providers"
	"github.com/ywl0806/yuno_kiroku/internal/routers"
	"github.com/ywl0806/yuno_kiroku/internal/services"
	"github.com/ywl0806/yuno_kiroku/internal/store"
	"github.com/ywl0806/yuno_kiroku/internal/validator"
)

// Initialize the root router on the app
func Init(e *echo.Echo) {
	// i18n 초기화
	i18n.Init()

	// database 초기화
	dbTx, err := sql.Open("postgres", viper.GetString("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	// 연결 테스트
	if err := dbTx.Ping(); err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// queries 초기화
	queries := db.New(dbTx)

	// store 초기화
	st := store.New(queries)

	storageProvider := providers.NewStorageProvider()
	// Storage service
	storageService := storageProvider.StorageService()

	// service (store 계층을 통해 데이터 접근)
	userService := services.NewUserService(st.User, st.Group)
	inviteService := services.NewInviteService(st.InviteToken, st.Group)
	authService := services.NewAuthService(
		userService,
		inviteService,
		viper.GetString("AUTH_SECRET_KEY"),
		&oauth.LineConfig{
			ChannelID:     viper.GetString("LINE_CHANNEL_ID"),
			ChannelSecret: viper.GetString("LINE_CHANNEL_SECRET"),
			CallbackURL:   viper.GetString("LINE_CALLBACK_URL"),
		},
		&oauth.KakaoConfig{
			ClientID:     viper.GetString("KAKAO_CLIENT_ID"),
			ClientSecret: viper.GetString("KAKAO_CLIENT_SECRET"),
			RedirectURI:  viper.GetString("KAKAO_REDIRECT_URI"),
		},
	)
	faceService := services.NewFaceService(st.Face, st.Identity, st.MediaItem)
	mediaItemService := services.NewMediaItemService(st.MediaItem, storageService, faceService)
	identityService := services.NewIdentityService(st.Identity)
	albumService := services.NewAlbumService(st.Album)

	// handler
	userHandler := handlers.NewUserHandler(userService)
	mediaItemHandler := handlers.NewMediaItemHandler(mediaItemService, faceService)
	authHandler := handlers.NewAuthHandler(authService)
	inviteHandler := handlers.NewInviteHandler(inviteService)
	identityHandler := handlers.NewIdentityHandler(identityService)
	albumHandler := handlers.NewAlbumHandler(albumService)

	// root router
	root := e.Group("/api")

	// health check
	root.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	// router 초기화
	userRouter := routers.NewUserRouter(*userHandler)
	mediaItemRouter := routers.NewMediaItemRouter(*mediaItemHandler)
	authRouter := routers.NewAuthRouter(*authHandler)
	inviteRouter := routers.NewInviteRouter(*inviteHandler)
	identityRouter := routers.NewIdentityRouter(*identityHandler)
	albumRouter := routers.NewAlbumRouter(*albumHandler)

	e.Validator = validator.NewCustomValidator()

	// 미들웨어 등록 (등록 순서대로 요청 시 실행됨 - RequestID가 먼저 와야 Logger에서 request_id 사용 가능)
	// 1. RequestID - 요청 ID 설정 (Logger보다 먼저 등록해야 Format에서 ${request_id} 출력됨)
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/api/health"
		},
	}))
	// 2. Recover - 패닉 복구
	e.Use(middleware.Recover())
	// 3. Logger - 로깅
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} uri=${uri},\n status=${status},\n latency=${latency_human}\n request_id=${id} \n\n",
	}))
	// 4. Locale - Accept-Language 기반 로케일 설정 (i18n)
	e.Use(middlewares.LocaleMiddleware())
	// 5. ErrorHandler - 에러 처리
	errorHandler := middlewares.NewErrorHandler()
	e.Use(errorHandler.Handler)
	// init routers
	userRouter.Register(root)
	mediaItemRouter.Register(root)
	authRouter.Register(root)
	inviteRouter.Register(root)
	identityRouter.Register(root)
	albumRouter.Register(root)

}
