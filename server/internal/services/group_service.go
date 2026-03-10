package services

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store"
)

type GroupService struct {
	familyStore store.FamilyStore
	groupStore  store.GroupStore
}

func NewGroupService(familyStore store.FamilyStore, groupStore store.GroupStore) *GroupService {
	return &GroupService{familyStore: familyStore, groupStore: groupStore}
}

func (s *GroupService) GetGroups(ctx context.Context, familyID int32) ([]db.Group, error) {
	return s.groupStore.FindGroupsByFamilyID(ctx, familyID)
}
