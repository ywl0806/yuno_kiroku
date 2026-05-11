package models

import (
	"time"

	"github.com/ywl0806/yuno_kiroku/internal/db"
	internalutils "github.com/ywl0806/yuno_kiroku/internal/utils"
)

type KidResponse struct {
	ID         int32      `json:"id"`
	FamilyID   string     `json:"family_id"`
	Name       *string    `json:"name"`
	BirthDate  *time.Time `json:"birth_date"`
	IdentityID *int32     `json:"identity_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type CreateKidRequest struct {
	Name       *string `json:"name"`
	BirthDate  *string `json:"birth_date"`
	IdentityID *int32  `json:"identity_id"`
}

type UpdateKidRequest struct {
	Name       *string `json:"name"`
	BirthDate  *string `json:"birth_date"`
	IdentityID *int32  `json:"identity_id"`
}

type GetKidsWithFaceImgRequest struct {
	Year  int `query:"year" validate:"required"`
	Month int `query:"month" validate:"required"`
}

type KidWithFaceImg struct {
	KidID        int32  `json:"kid_id"`
	Name         string `json:"name"`
	FaceImgURL   string `json:"face_img_url"`
	MediaItemID  string `json:"media_item_id"`
	TakenAtYear  int    `json:"taken_at_year"`
	TakenAtMonth int    `json:"taken_at_month"`
	BirthDate    string `json:"birth_date"`
}

func NewKidWithFaceImg(k *db.GetKidsWithRandomFaceImgRow) *KidWithFaceImg {
	return &KidWithFaceImg{
		KidID:        k.ID,
		Name:         k.Name.String,
		FaceImgURL:   internalutils.ParseStoragePath(k.StorageKey),
		MediaItemID:  k.MediaItemID,
		TakenAtYear:  k.TakenAt.Year(),
		TakenAtMonth: int(k.TakenAt.Month()),
		BirthDate:    k.BirthDate.Time.Format("2006-01-02"),
	}
}

type KidWithFaceImgResponse struct {
	Kids     []KidWithFaceImg `json:"kids"`
	FastKids []KidWithFaceImg `json:"fast_kids"`
}

func NewKidWithFaceImgResponse(ks []db.GetKidsWithRandomFaceImgRow, fastKs []db.GetKidsWithRandomFaceImgRow) *KidWithFaceImgResponse {

	kids := make([]KidWithFaceImg, len(ks))
	for i, k := range ks {
		kids[i] = *NewKidWithFaceImg(&k)
	}

	fastKids := make([]KidWithFaceImg, len(fastKs))
	for i, k := range fastKs {
		fastKids[i] = *NewKidWithFaceImg(&k)
	}
	return &KidWithFaceImgResponse{
		Kids:     kids,
		FastKids: fastKids,
	}
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
