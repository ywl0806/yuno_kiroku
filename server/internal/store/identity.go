package store

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

// IdentityStore identity(인물) 데이터 접근 인터페이스
type IdentityStore interface {
	CreateIdentity(ctx context.Context, arg db.CreateIdentityParams) (db.Identity, error)
	FindIdentityByIdAndGroupId(ctx context.Context, arg db.FindIdentityByIdAndGroupIdParams) (db.Identity, error)
	FindIdentitiesByGroupId(ctx context.Context, groupID int32) ([]db.Identity, error)
	UpdateIdentityByIdAndGroupId(ctx context.Context, arg db.UpdateIdentityByIdAndGroupIdParams) (db.Identity, error)
}

type identityStore struct {
	queries *db.Queries
}

// NewIdentityStore IdentityStore 구현체 생성
func NewIdentityStore(queries *db.Queries) IdentityStore {
	return &identityStore{queries: queries}
}

func (s *identityStore) CreateIdentity(ctx context.Context, arg db.CreateIdentityParams) (db.Identity, error) {
	return wrapErr(s.queries.CreateIdentity(ctx, arg))
}

func (s *identityStore) FindIdentityByIdAndGroupId(ctx context.Context, arg db.FindIdentityByIdAndGroupIdParams) (db.Identity, error) {
	identity, err := s.queries.FindIdentityByIdAndGroupId(ctx, arg)
	return wrapErr(identity, err, "field.identity")
}

func (s *identityStore) FindIdentitiesByGroupId(ctx context.Context, groupID int32) ([]db.Identity, error) {
	return wrapErr(s.queries.FindIdentitiesByGroupId(ctx, groupID))
}

func (s *identityStore) UpdateIdentityByIdAndGroupId(ctx context.Context, arg db.UpdateIdentityByIdAndGroupIdParams) (db.Identity, error) {
	return wrapErr(s.queries.UpdateIdentityByIdAndGroupId(ctx, arg))
}
