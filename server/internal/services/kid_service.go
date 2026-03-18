package services

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store"
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
