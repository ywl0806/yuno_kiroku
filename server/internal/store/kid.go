package store

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

// KidStore 아이 데이터 접근 인터페이스
type KidStore interface {
	FindKidsByFamilyID(ctx context.Context, familyID int32) ([]db.Kid, error)
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
