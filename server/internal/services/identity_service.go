package services

import (
	"context"
	"database/sql"

	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store"
)

type IdentityService struct {
	identityStore store.IdentityStore
}

func NewIdentityService(identityStore store.IdentityStore) *IdentityService {
	return &IdentityService{identityStore: identityStore}
}

func (s *IdentityService) CreateIdentity(ctx context.Context, name string, groupId int32) (db.Identity, error) {
	identity, err := s.identityStore.CreateIdentity(ctx, db.CreateIdentityParams{
		Name:    sql.NullString{String: name, Valid: true},
		GroupID: groupId,
	})
	if err != nil {
		return db.Identity{}, err
	}
	return identity, nil
}

func (s *IdentityService) FindIdentityByIdAndGroupId(ctx context.Context, id int32, groupId int32) (db.Identity, error) {
	identity, err := s.identityStore.FindIdentityByIdAndGroupId(ctx, db.FindIdentityByIdAndGroupIdParams{
		ID:      id,
		GroupID: groupId,
	})
	if err != nil {
		return db.Identity{}, err
	}
	return identity, nil
}

func (s *IdentityService) FindIdentitiesByGroupId(ctx context.Context, groupId int32) ([]db.Identity, error) {
	identities, err := s.identityStore.FindIdentitiesByGroupId(ctx, groupId)
	if err != nil {
		return nil, err
	}
	if identities == nil {
		return []db.Identity{}, nil
	}
	return identities, nil
}

func (s *IdentityService) UpdateIdentityByIdAndGroupId(ctx context.Context, id int32, groupId int32, name string) (db.Identity, error) {

	identity, err := s.identityStore.UpdateIdentityByIdAndGroupId(ctx, db.UpdateIdentityByIdAndGroupIdParams{
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
	if err != nil {
		return err
	}
	return nil
}
