package services

import "github.com/ywl0806/yuno_kiroku/internal/store"

type GroupService struct {
	groupStore     store.GroupStore
	clanGroupStore store.ClanGroupStore
}

func NewGroupService(groupStore store.GroupStore) *GroupService {
	return &GroupService{groupStore: groupStore}
}
