package services

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store"
)

type IdentityService struct {
	identityStore        store.IdentityStore
	identityFaceImgStore store.IdentityFaceImgStore
}

func NewIdentityService(identityStore store.IdentityStore, identityFaceImgStore store.IdentityFaceImgStore) *IdentityService {
	return &IdentityService{identityStore: identityStore, identityFaceImgStore: identityFaceImgStore}
}

func (s *IdentityService) CreateIdentity(ctx context.Context, familyId string) (db.Identity, error) {
	identity, err := s.identityStore.CreateIdentity(ctx, familyId)
	if err != nil {
		return db.Identity{}, err
	}
	return identity, nil
}

func (s *IdentityService) FindIdentityByIdAndFamilyId(ctx context.Context, id int32, familyId string) (db.Identity, error) {
	identity, err := s.identityStore.FindIdentityByIdAndFamilyId(ctx, db.FindIdentityByIdAndFamilyIdParams{
		ID:       id,
		FamilyID: familyId,
	})
	if err != nil {
		return db.Identity{}, err
	}
	return identity, nil
}

func (s *IdentityService) FindIdentitiesByFamilyId(ctx context.Context, familyId string) ([]db.Identity, error) {
	identities, err := s.identityStore.FindIdentitiesByFamilyId(ctx, familyId)
	if err != nil {
		return nil, err
	}
	if identities == nil {
		return []db.Identity{}, nil
	}
	return identities, nil
}

// 가장 최근의 identity 얼굴 이미지 조회
func (s *IdentityService) FindNewestIdentityFaceImgByFamilyId(ctx context.Context, familyId string, onlyNotLinked bool, withKidIds []int32, withUserIds []string) ([]db.FindNewestIdentityFaceImgByFamilyIdRow, error) {
	identityFaceImgs, err := s.identityFaceImgStore.FindNewestIdentityFaceImgByFamilyId(ctx, db.FindNewestIdentityFaceImgByFamilyIdParams{
		FamilyID:      familyId,
		OnlyNotLinked: onlyNotLinked,
		WithKidIds:    withKidIds,
		WithUserIds:   withUserIds,
	})
	if err != nil {
		return nil, err
	}
	return identityFaceImgs, nil
}

// identity 옵션 조회
func (s *IdentityService) GetIdentityOptions(ctx context.Context, familyID string) ([]db.GetIdentityOptionsRow, error) {
	return s.identityStore.GetIdentityOptions(ctx, familyID)
}
