package models

import (
	"time"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

type KidResponse struct {
	ID         int32      `json:"id"`
	FamilyID   int32      `json:"family_id"`
	Name       *string    `json:"name"`
	BirthDate  *time.Time `json:"birth_date"`
	IdentityID *int32     `json:"identity_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func NewKidResponse(k *db.Kid) *KidResponse {
	var name *string
	if k.Name.Valid {
		name = &k.Name.String
	}

	var birthDate *time.Time
	if k.BirthDate.Valid {
		birthDate = &k.BirthDate.Time
	}

	var identityID *int32
	if k.IdentityID.Valid {
		identityID = &k.IdentityID.Int32
	}

	return &KidResponse{
		ID:         k.ID,
		FamilyID:   k.FamilyID,
		Name:       name,
		BirthDate:  birthDate,
		IdentityID: identityID,
		CreatedAt:  k.CreatedAt,
		UpdatedAt:  k.UpdatedAt,
	}
}
