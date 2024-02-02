package user

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/api/user/store"
)

type UserController struct {
	userStore store.UserStore
}

func NewUserController(userStore store.UserStore) *UserController {
	return &UserController{userStore: userStore}
}

func (h *UserController) GetUser(c echo.Context) error {
	log.Default()
	h.userStore.CreateUser()
	return nil
}
