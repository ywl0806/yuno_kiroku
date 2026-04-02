package api

import (
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"

	"github.com/labstack/echo/v4/middleware"
	"github.com/ywl0806/yuno_kiroku/internal/api/handlers"
	"github.com/ywl0806/yuno_kiroku/internal/api/middlewares"
	"github.com/ywl0806/yuno_kiroku/internal/api/routers"
	"github.com/ywl0806/yuno_kiroku/internal/consts"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/i18n"
	"github.com/ywl0806/yuno_kiroku/internal/oauth"
	"github.com/ywl0806/yuno_kiroku/internal/providers"
	"github.com/ywl0806/yuno_kiroku/internal/services"
	"github.com/ywl0806/yuno_kiroku/internal/store"
	"github.com/ywl0806/yuno_kiroku/internal/validator"
)

// Initialize the root router on the app
func Init(e *echo.Echo) {
	// i18n 초기화
	i18n.Init()

	// DB 초기화
	sqlDB, dbTx, err := db.Init(viper.GetString("APP_MODE") == "dev" || viper.GetString("APP_MODE") == "local" || viper.GetString("APP_MODE") == "local_dev")
	if err != nil {
		log.Fatalf("Failed to initialize DB: %v", err)
	}

	// queries 초기화
	queries := db.New(dbTx)

	// store 초기화
	st := store.New(sqlDB, queries)

	storageProvider := providers.NewStorageProvider()
	// Storage service
	storageService := storageProvider.StorageService()

	// service (store 계층을 통해 데이터 접근)
	userService := services.NewUserService(st.User, st.Family, st.Group)
	inviteService := services.NewInviteService(st.InviteToken, st.Family, st.Group)
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
	imageUploader := services.NewImageUploader(storageService)
	mediaItemService := services.NewMediaItemService(st.MediaItem, imageUploader)
	identityService := services.NewIdentityService(st.Identity, st.IdentityFaceImg)
	albumService := services.NewAlbumService(st.Album, st.AlbumGroupPermission, st)
	groupService := services.NewGroupService(st.Family, st.Group)
	kidService := services.NewKidService(st.Kid)

	// handler
	userHandler := handlers.NewUserHandler(userService)
	mediaItemHandler := handlers.NewMediaItemHandler(mediaItemService)
	authHandler := handlers.NewAuthHandler(authService)
	inviteHandler := handlers.NewInviteHandler(inviteService)
	identityHandler := handlers.NewIdentityHandler(identityService)
	albumHandler := handlers.NewAlbumHandler(albumService)
	groupHandler := handlers.NewGroupHandler(groupService)
	kidHandler := handlers.NewKidHandler(kidService)
	settingsHandler := handlers.NewSettingsHandler(groupService, userService, albumService, kidService, identityService)

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
	groupRouter := routers.NewGroupRouter(*groupHandler)
	kidRouter := routers.NewKidRouter(*kidHandler)
	settingsRouter := routers.NewSettingsRouter(*settingsHandler)

	e.Validator = validator.NewCustomValidator()

	// 미들웨어 등록 (등록 순서대로 요청 시 실행됨 - RequestID가 먼저 와야 Logger에서 request_id 사용 가능)
	// 1. RequestID - 요청 ID 설정 (Logger보다 먼저 등록해야 Format에서 ${request_id} 출력됨)
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/api/health"
		},
		RequestIDHandler: func(c echo.Context, requestID string) {
			req := c.Request()
			ctx := context.WithValue(req.Context(), consts.RequestIDKey, requestID)
			c.SetRequest(req.WithContext(ctx))
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
	groupRouter.Register(root)
	kidRouter.Register(root)
	settingsRouter.Register(root)

}
