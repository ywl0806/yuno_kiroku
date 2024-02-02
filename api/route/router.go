package route

import (
	"github.com/labstack/echo/v4"

	"github.com/ywl0806/yuno_kiroku/api/db"
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

	// handler
	userHandler := user.NewUserContoller(mUserStore)

	// root router
	root := e.Group("/api")

	// init routers
	user.Register(root, *userHandler)
}
