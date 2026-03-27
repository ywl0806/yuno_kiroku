package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store"
	"github.com/ywl0806/yuno_kiroku/pkg/utils"
)

type KidService struct {
	kidStore store.KidStore
}

func NewKidService(kidStore store.KidStore) *KidService {
	return &KidService{kidStore: kidStore}
}

func (s *KidService) GetKids(ctx context.Context, familyID int32) ([]db.Kid, error) {
	return s.kidStore.FindKidsByFamilyID(ctx, familyID)
}

func (s *KidService) GetKidByID(ctx context.Context, id int32) (db.Kid, error) {
	return s.kidStore.GetKidByID(ctx, id)
}

func (s *KidService) CreateKid(ctx context.Context, familyID int32, name *string, birthDate *time.Time, identityID *int32) (db.Kid, error) {
	params := db.CreateKidParams{FamilyID: familyID}
	if name != nil {
		params.Name = sql.NullString{String: *name, Valid: true}
	}
	if birthDate != nil {
		params.BirthDate = sql.NullTime{Time: *birthDate, Valid: true}
	}
	if identityID != nil {
		params.IdentityID = sql.NullInt32{Int32: *identityID, Valid: true}
	}
	return s.kidStore.CreateKid(ctx, params)
}

func (s *KidService) UpdateKid(ctx context.Context, id int32, name *string, birthDate *time.Time, identityID *int32) (db.Kid, error) {
	params := db.UpdateKidParams{ID: id}
	if name != nil {
		params.Name = sql.NullString{String: *name, Valid: true}
	}
	if birthDate != nil {
		params.BirthDate = sql.NullTime{Time: *birthDate, Valid: true}
	}
	if identityID != nil {
		params.IdentityID = sql.NullInt32{Int32: *identityID, Valid: true}
	}
	return s.kidStore.UpdateKid(ctx, params)
}

func (s *KidService) DeleteKid(ctx context.Context, id int32) error {
	return s.kidStore.DeleteKid(ctx, id)
}

type KidWithFaceImg struct {
	KidID       int32
	Name        string
	FaceImgURL  string
	MediaItemID int32
	TakenAt     time.Time
}

// GetKidsWithFaceImg 월별 아이 얼굴 사진 조회
func (s *KidService) GetKidsWithFaceImg(ctx context.Context, familyID int32, year int, month int) ([]db.GetKidsWithRandomFaceImgRow, error) {
	return s.kidStore.GetKidsWithRandomFaceImg(ctx, db.GetKidsWithRandomFaceImgParams{
		FamilyID:    familyID,
		TakenAtTo:   utils.GetLastDayOfMonth(year, month),
		TakenAtFrom: utils.GetFirstDayOfMonth(year, month),
	})
}
