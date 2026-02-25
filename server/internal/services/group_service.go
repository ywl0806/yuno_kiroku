package services

import "github.com/ywl0806/yuno_kiroku/internal/store"

type GroupService struct {
	familyStore store.FamilyStore
	groupStore  store.GroupStore
}

func NewGroupService(familyStore store.FamilyStore, groupStore store.GroupStore) *GroupService {
	return &GroupService{familyStore: familyStore, groupStore: groupStore}
}
