package route

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/api/user"
)

func New(e *echo.Echo) {
	fmt.Println()
	g := e.Group("/api")
	user.Route(g)

}
