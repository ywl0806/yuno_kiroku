package user

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/api/user/store"
)

type UserContoller struct {
	userStore store.UserStore
}

func NewUserContoller(userStore store.UserStore) *UserContoller {
	return &UserContoller{userStore: userStore}
}

func (h *UserContoller) GetUser(c echo.Context) error {
	log.Default()
	h.userStore.CreateUser()
	return nil
}
