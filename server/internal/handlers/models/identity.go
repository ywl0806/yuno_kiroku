package models

import (
	"time"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

type IdentityResponse struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	FamilyId  int32     `json:"family_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewIdentityResponse(identity *db.Identity) *IdentityResponse {
	return &IdentityResponse{
		ID:        identity.ID,
		Name:      identity.Name.String,
		FamilyId:  identity.FamilyID,
		CreatedAt: identity.CreatedAt,
		UpdatedAt: identity.UpdatedAt,
	}
}

type IdentityListResponse struct {
	Identities []IdentityResponse `json:"identities"`
}

func NewIdentityListResponse(identities []db.Identity) *IdentityListResponse {
	identityResponses := make([]IdentityResponse, len(identities))
	for i, identity := range identities {
		identityResponses[i] = *NewIdentityResponse(&identity)
	}
	return &IdentityListResponse{
		Identities: identityResponses,
	}
}

type CreateIdentityRequest struct {
	Name string `json:"name"`
}

type FindIdentityByIdAndGroupIdRequest struct {
	ID int32 `json:"id"`
}

type UpdateIdentityByIdAndGroupIdRequest struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}
