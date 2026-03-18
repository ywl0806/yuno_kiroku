package services

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store"
)

type IdentityService struct {
	identityStore store.IdentityStore
}

func NewIdentityService(identityStore store.IdentityStore) *IdentityService {
	return &IdentityService{identityStore: identityStore}
}

func (s *IdentityService) CreateIdentity(ctx context.Context, familyId int32) (db.Identity, error) {
	identity, err := s.identityStore.CreateIdentity(ctx, familyId)
	if err != nil {
		return db.Identity{}, err
	}
	return identity, nil
}

func (s *IdentityService) FindIdentityByIdAndFamilyId(ctx context.Context, id int32, familyId int32) (db.Identity, error) {
	identity, err := s.identityStore.FindIdentityByIdAndFamilyId(ctx, db.FindIdentityByIdAndFamilyIdParams{
		ID:       id,
		FamilyID: familyId,
	})
	if err != nil {
		return db.Identity{}, err
	}
	return identity, nil
}

func (s *IdentityService) FindIdentitiesByFamilyId(ctx context.Context, familyId int32) ([]db.Identity, error) {
	identities, err := s.identityStore.FindIdentitiesByFamilyId(ctx, familyId)
	if err != nil {
		return nil, err
	}
	if identities == nil {
		return []db.Identity{}, nil
	}
	return identities, nil
}
