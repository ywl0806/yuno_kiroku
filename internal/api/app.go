package api

import (
	"github.com/labstack/echo/v4"

	"github.com/ywl0806/yuno_kiroku/pkg/db"
	"github.com/ywl0806/yuno_kiroku/pkg/storage"

	"github.com/ywl0806/yuno_kiroku/internal/api/services/auth"
	"github.com/ywl0806/yuno_kiroku/internal/api/services/photo"
	photoStore "github.com/ywl0806/yuno_kiroku/internal/api/services/photo/store"
	"github.com/ywl0806/yuno_kiroku/internal/api/services/user"
	userStore "github.com/ywl0806/yuno_kiroku/internal/api/services/user/store"
)

// Initialize the root router on the app
func Init(e *echo.Echo) {

	// db
	client := db.ConnectDB()
	db := client.Database("yuno")

	// store
	userStore := userStore.NewUserStore(db)
	photoStore := photoStore.NewPhotoStore(db)

	// Storage service
	sStorage := storage.NewLocalStorageService("standard")
	lStorage := storage.NewLocalStorageService("longterm")

	// Controller
	userController := user.NewUserController(userStore)
	photoController := photo.NewPhotoController(photoStore, sStorage, lStorage)

	// root router
	root := e.Group("/api")

	// health check
	root.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})
	// init routers
	user.Register(root, *userController)
	photo.Register(root, *photoController)
	// auth
	authController := auth.NewAuthController(userStore)
	auth.Register(root, *authController)
}
