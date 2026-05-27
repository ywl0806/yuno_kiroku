package models

import (
	"time"

	"github.com/ywl0806/yuno_kiroku/internal/db"
	internalutils "github.com/ywl0806/yuno_kiroku/internal/utils"
)

type IdentityResponse struct {
	ID        int32     `json:"id"`
	FamilyId  string    `json:"family_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewIdentityResponse(identity *db.Identity) *IdentityResponse {
	return &IdentityResponse{
		ID:        identity.ID,
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

type FindIdentityByIdAndGroupIdRequest struct {
	ID int32 `json:"id"`
}

type IdentityOptionResponse struct {
	ID       int32   `json:"id"`
	KidID    *int32  `json:"kid_id"`
	KidName  *string `json:"kid_name"`
	UserID   *string `json:"user_id"`
	UserName *string `json:"user_name"`
	ImageURL *string `json:"image_url"`
}

func NewIdentityOptionResponse(identityOption *db.GetIdentityOptionsRow) *IdentityOptionResponse {
	var kidID *int32
	if identityOption.KidID.Valid {
		kidID = &identityOption.KidID.Int32
	}
	var userID *string
	if identityOption.UserID.Valid {
		v := identityOption.UserID.UUID.String()
		userID = &v
	}
	var kidName *string
	if identityOption.KidName.Valid {
		kidName = &identityOption.KidName.String
	}
	var userName *string
	if identityOption.UserName.Valid {
		userName = &identityOption.UserName.String
	}
	var imageURL *string
	if identityOption.StorageKey.Valid {
		url := internalutils.ParseStoragePath(identityOption.StorageKey.String)
		imageURL = &url
	}
	return &IdentityOptionResponse{
		ID:       identityOption.ID,
		KidID:    kidID,
		KidName:  kidName,
		UserID:   userID,
		UserName: userName,
		ImageURL: imageURL,
	}
}

func NewIdentityOptionResponses(identityOptions []db.GetIdentityOptionsRow) *[]IdentityOptionResponse {
	identityOptionResponses := make([]IdentityOptionResponse, len(identityOptions))
	for i, identityOption := range identityOptions {
		identityOptionResponses[i] = *NewIdentityOptionResponse(&identityOption)
	}
	return &identityOptionResponses
}
