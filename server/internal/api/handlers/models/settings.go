package models

import (
	"github.com/ywl0806/yuno_kiroku/internal/db"
	internalutils "github.com/ywl0806/yuno_kiroku/internal/utils"
)

type SettingsDataResponse struct {
	Groups  []GroupResponse  `json:"groups"`
	Members []MemberResponse `json:"members"`
	Albums  []AlbumResponse  `json:"albums"`
	Kids    []KidResponse    `json:"kids"`
}

type IdentityFaceOptionResponse struct {
	IdentityID int32  `json:"identity_id"`
	ImageURL   string `json:"image_url"`
}

func NewIdentityFaceOptionResponses(identityFaceOptions []db.FindNewestIdentityFaceImgByFamilyIdRow) *[]IdentityFaceOptionResponse {
	identityFaceOptionResponses := make([]IdentityFaceOptionResponse, len(identityFaceOptions))
	for i, identityFaceOption := range identityFaceOptions {
		identityFaceOptionResponses[i] = *NewIdentityFaceOptionResponse(&identityFaceOption)
	}
	return &identityFaceOptionResponses
}

func NewIdentityFaceOptionResponse(identityFaceOption *db.FindNewestIdentityFaceImgByFamilyIdRow) *IdentityFaceOptionResponse {
	return &IdentityFaceOptionResponse{
		IdentityID: identityFaceOption.IdentityID,
		ImageURL:   internalutils.ParseStoragePath(identityFaceOption.StorageKey),
	}
}

type IdentityFaceOptionRequest struct {
	OnlyNotLinked bool    `query:"only_not_linked"`
	WithKidIds    []int32 `query:"with_kid_ids"`
	WithUserIds   []int32 `query:"with_user_ids"`
}
