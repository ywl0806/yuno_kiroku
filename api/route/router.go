package route

import (
	"github.com/labstack/echo/v4"

	"github.com/ywl0806/yuno_kiroku/api/db"
	"github.com/ywl0806/yuno_kiroku/api/user"
)

func Init(e *echo.Echo) {

	client := db.New()
	root := e.Group("/api")
	user.Register(root, client)

}
