package models

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/pkg/utils"
)

type CreateUserRequest struct {
	Name        string `json:"name" validate:"required"`
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required,min=6"`
	GroupID     int32  `json:"group_id" validate:"required"`
	ClanGroupID int32  `json:"clan_group_id" validate:"required"`
}

func (CreateUserRequest) Bind(c echo.Context, params *db.CreateUserParams) error {
	req := new(CreateUserRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}
	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}
	params.Name = sql.NullString{String: req.Name, Valid: true}
	params.Username = req.Username
	params.Password = hashPassword
	params.GroupID = req.GroupID
	params.ClanGroupID = req.ClanGroupID
	return nil
}
