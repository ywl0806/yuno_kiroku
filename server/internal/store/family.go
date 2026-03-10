package store

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

// FamilyStore 가족(family) 데이터 접근 인터페이스
type FamilyStore interface {
	FindFamilyByID(ctx context.Context, id int32) (db.FindFamilyByIDRow, error)
	GetFamilyByID(ctx context.Context, id int32) (db.Family, error)
	CreateFamily(ctx context.Context, name string) (db.Family, error)
}

type familyStore struct {
	queries *db.Queries
}

// NewFamilyStore FamilyStore 구현체 생성
func NewFamilyStore(queries *db.Queries) FamilyStore {
	return &familyStore{queries: queries}
}

func (s *familyStore) FindFamilyByID(ctx context.Context, id int32) (db.FindFamilyByIDRow, error) {
	row, err := s.queries.FindFamilyByID(ctx, id)
	return wrapErr(row, err, "field.family")
}

func (s *familyStore) GetFamilyByID(ctx context.Context, id int32) (db.Family, error) {
	return wrapErr(s.queries.GetFamilyByID(ctx, id))
}

func (s *familyStore) CreateFamily(ctx context.Context, name string) (db.Family, error) {
	return wrapErr(s.queries.CreateFamily(ctx, name))
}
