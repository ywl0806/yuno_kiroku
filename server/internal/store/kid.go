package store

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

// KidStore 아이 데이터 접근 인터페이스
type KidStore interface {
	FindKidsByFamilyID(ctx context.Context, familyID int32) ([]db.Kid, error)
	CreateKid(ctx context.Context, arg db.CreateKidParams) (db.Kid, error)
	UpdateKid(ctx context.Context, arg db.UpdateKidParams) (db.Kid, error)
	DeleteKid(ctx context.Context, id int32) error
	GetKidByID(ctx context.Context, id int32) (db.Kid, error)
	GetKidsWithRandomFaceImg(ctx context.Context, arg db.GetKidsWithRandomFaceImgParams) ([]db.GetKidsWithRandomFaceImgRow, error)
}

type kidStore struct {
	queries *db.Queries
}

// NewKidStore KidStore 구현체 생성
func NewKidStore(queries *db.Queries) KidStore {
	return &kidStore{queries: queries}
}

func (s *kidStore) FindKidsByFamilyID(ctx context.Context, familyID int32) ([]db.Kid, error) {
	return wrapErr(s.queries.FindKidsByFamilyID(ctx, familyID))
}

func (s *kidStore) CreateKid(ctx context.Context, arg db.CreateKidParams) (db.Kid, error) {
	return wrapErr(s.queries.CreateKid(ctx, arg))
}

func (s *kidStore) UpdateKid(ctx context.Context, arg db.UpdateKidParams) (db.Kid, error) {
	return wrapErr(s.queries.UpdateKid(ctx, arg))
}

func (s *kidStore) DeleteKid(ctx context.Context, id int32) error {
	return mapDBError(s.queries.DeleteKid(ctx, id))
}

func (s *kidStore) GetKidByID(ctx context.Context, id int32) (db.Kid, error) {
	return wrapErr(s.queries.GetKidByID(ctx, id))
}

func (s *kidStore) GetKidsWithRandomFaceImg(ctx context.Context, arg db.GetKidsWithRandomFaceImgParams) ([]db.GetKidsWithRandomFaceImgRow, error) {
	return wrapErr(s.queries.GetKidsWithRandomFaceImg(ctx, arg))
}
