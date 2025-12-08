package services

import (
	"context"
	"database/sql"

	apiErrors "github.com/ywl0806/yuno_kiroku/internal/api/errors"
	"github.com/ywl0806/yuno_kiroku/internal/db"
)

type IdentityService struct {
	queries *db.Queries
}

func NewIdentityService(queries *db.Queries) *IdentityService {
	return &IdentityService{queries: queries}
}

func (s *IdentityService) CreateIdentity(ctx context.Context, name string, groupId int32) (db.Identity, error) {
	identity, err := s.queries.CreateIdentity(ctx, db.CreateIdentityParams{
		Name:    sql.NullString{String: name, Valid: true},
		GroupID: groupId,
	})
	if err != nil {
		return db.Identity{}, err
	}
	return identity, nil
}

func (s *IdentityService) FindIdentityByIdAndGroupId(ctx context.Context, id int32, groupId int32) (db.Identity, error) {
	identity, err := s.queries.FindIdentityByIdAndGroupId(ctx, db.FindIdentityByIdAndGroupIdParams{
		ID:      id,
		GroupID: groupId,
	})
	if err == sql.ErrNoRows {
		return db.Identity{}, apiErrors.ErrIdentityNotFound
	}
	if err != nil {
		return db.Identity{}, err
	}
	return identity, nil
}

func (s *IdentityService) FindIdentitiesByGroupId(ctx context.Context, groupId int32) ([]db.Identity, error) {
	identities, err := s.queries.FindIdentitiesByGroupId(ctx, groupId)
	if err != nil {
		return nil, err
	}
	// 빈 배열은 정상적인 결과이므로 에러가 아님
	if identities == nil {
		return []db.Identity{}, nil
	}
	return identities, nil
}

func (s *IdentityService) UpdateIdentityByIdAndGroupId(ctx context.Context, id int32, groupId int32, name string) (db.Identity, error) {

	identity, err := s.queries.UpdateIdentityByIdAndGroupId(ctx, db.UpdateIdentityByIdAndGroupIdParams{
		ID:      id,
		GroupID: groupId,
		Name:    sql.NullString{String: name, Valid: true},
	})
	if err != nil {
		return db.Identity{}, err
	}
	return identity, nil
}

func (s *IdentityService) ValidateUpdateIdentityByIdAndGroupId(ctx context.Context, id int32, groupId int32, name string) error {
	_, err := s.FindIdentityByIdAndGroupId(ctx, id, groupId)
	if err == apiErrors.ErrIdentityNotFound {
		return apiErrors.ErrIdentityNotFound
	}
	if err != nil {
		return err
	}
	return nil
}
