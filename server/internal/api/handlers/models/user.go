package models

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/pkg/utils"
)

type MeResponse struct {
	ID       int32   `json:"id"`
	Name     *string `json:"name"`
	Username string  `json:"username"`
	Provider *string `json:"provider"`
}

func NewMeResponse(u *db.User) *MeResponse {
	var name *string
	if u.Name.Valid {
		name = &u.Name.String
	}
	var provider *string
	if u.Provider.Valid {
		provider = &u.Provider.String
	}
	return &MeResponse{
		ID:       u.ID,
		Name:     name,
		Username: u.Username,
		Provider: provider,
	}
}

type UpdateMeRequest struct {
	Name string `json:"name"`
}

type UpdateMemberRequest struct {
	GroupID           int32  `json:"group_id"`
	FamilyTitle       string `json:"family_title"`
	CustomFamilyTitle string `json:"custom_family_title"`
}

type MemberResponse struct {
	ID                int32   `json:"id"`
	Name              string  `json:"name"`
	Username          string  `json:"username"`
	FamilyID          int32   `json:"family_id"`
	GroupID           int32   `json:"group_id"`
	FamilyTitle       *string `json:"family_title"`
	CustomFamilyTitle *string `json:"custom_family_title"`
}

func NewMemberResponse(user *db.User) *MemberResponse {
	var familyTitle *string
	if user.FamilyTitle.Valid {
		familyTitle = &user.FamilyTitle.String
	}
	var customFamilyTitle *string
	if user.CustomFamilyTitle.Valid {
		customFamilyTitle = &user.CustomFamilyTitle.String
	}
	return &MemberResponse{
		ID:                user.ID,
		Name:              user.Name.String,
		Username:          user.Username,
		FamilyID:          user.FamilyID,
		GroupID:           user.GroupID,
		FamilyTitle:       familyTitle,
		CustomFamilyTitle: customFamilyTitle,
	}
}

type UpdateMemberFamilyTitleRequest struct {
	FamilyTitle       string `json:"family_title"`
	CustomFamilyTitle string `json:"custom_family_title"`
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
	FamilyID int32  `json:"family_id" validate:"required"`
	GroupID  int32  `json:"group_id" validate:"required"`
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
	params.FamilyID = req.FamilyID
	params.GroupID = req.GroupID
	return nil
}
