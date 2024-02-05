package api

import (
	"github.com/labstack/echo/v4"

	"github.com/ywl0806/yuno_kiroku/api/db"
	"github.com/ywl0806/yuno_kiroku/api/lib/storage"
	"github.com/ywl0806/yuno_kiroku/api/photo"
	photoStore "github.com/ywl0806/yuno_kiroku/api/photo/store"
	"github.com/ywl0806/yuno_kiroku/api/user"
	userStore "github.com/ywl0806/yuno_kiroku/api/user/store"
)

// Initialize the root router on the app
func Init(e *echo.Echo) {

	// db
	client := db.ConnectDB()
	db := client.Database("yuno")

	// store
	mUserStore := userStore.NewMUserStore(db)
	mPhotoStore := photoStore.NewMPhotoStore(db)

	// Storage service
	sStorage := storage.NewLocalStorageService("standard")
	lStorage := storage.NewLocalStorageService("longterm")

	// Controller
	userController := user.NewUserController(mUserStore)
	photoController := photo.NewPhotoController(mPhotoStore, sStorage, lStorage)

	// root router
	root := e.Group("/api")

	// init routers
	user.Register(root, *userController)
	photo.Register(root, *photoController)
}
