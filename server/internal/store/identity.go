package store

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

// IdentityStore identity(인물) 데이터 접근 인터페이스
type IdentityStore interface {
	CreateIdentity(ctx context.Context, familyID int32) (db.Identity, error)
	FindIdentityByIdAndFamilyId(ctx context.Context, arg db.FindIdentityByIdAndFamilyIdParams) (db.Identity, error)
	FindIdentitiesByFamilyId(ctx context.Context, familyID int32) ([]db.Identity, error)
	GetIdentityOptions(ctx context.Context, familyID int32) ([]db.GetIdentityOptionsRow, error)
}

type identityStore struct {
	queries *db.Queries
}

// NewIdentityStore IdentityStore 구현체 생성
func NewIdentityStore(queries *db.Queries) IdentityStore {
	return &identityStore{queries: queries}
}

func (s *identityStore) CreateIdentity(ctx context.Context, familyID int32) (db.Identity, error) {
	return wrapErr(s.queries.CreateIdentity(ctx, familyID))
}

func (s *identityStore) FindIdentityByIdAndFamilyId(ctx context.Context, arg db.FindIdentityByIdAndFamilyIdParams) (db.Identity, error) {
	identity, err := s.queries.FindIdentityByIdAndFamilyId(ctx, arg)
	return wrapErr(identity, err, "field.identity")
}

func (s *identityStore) FindIdentitiesByFamilyId(ctx context.Context, familyID int32) ([]db.Identity, error) {
	return wrapErr(s.queries.FindIdentitiesByFamilyId(ctx, familyID))
}

func (s *identityStore) GetIdentityOptions(ctx context.Context, familyID int32) ([]db.GetIdentityOptionsRow, error) {
	return wrapErr(s.queries.GetIdentityOptions(ctx, familyID))
}
