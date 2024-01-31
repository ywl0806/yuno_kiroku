package user

import (
	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/api/user/store"
)

func Register(root *echo.Group, db *firestore.Client) {
	g := root.Group("/user")
	store := store.NewUserStore(db)
	g.GET("/hi", func(c echo.Context) error {

		store.Add()
		return c.String(200, "hello")
	})

}
