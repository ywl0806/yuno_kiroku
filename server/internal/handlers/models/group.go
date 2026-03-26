package models

import (
	"time"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

type GroupResponse struct {
	ID        int32     `json:"id"`
	FamilyID  int32     `json:"family_id"`
	IsAdmin   bool      `json:"is_admin"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewGroupResponse(g *db.Group) *GroupResponse {
	return &GroupResponse{
		ID:        g.ID,
		FamilyID:  g.FamilyID,
		IsAdmin:   g.IsAdmin,
		Name:      g.Name,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
}

type CreateGroupRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateGroupRequest struct {
	Name string `json:"name" validate:"required"`
}
