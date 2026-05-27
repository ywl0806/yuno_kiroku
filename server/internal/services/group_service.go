package services

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/apperr"
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

func (s *GroupService) GetGroups(ctx context.Context, familyID string) ([]db.Group, error) {
	return s.groupStore.FindGroupsByFamilyID(ctx, familyID)
}

func (s *GroupService) CreateGroup(ctx context.Context, familyID string, name string) (db.Group, error) {
	return s.groupStore.CreateGroup(ctx, db.CreateGroupParams{
		FamilyID: familyID,
		IsAdmin:  false,
		Name:     name,
	})
}

func (s *GroupService) UpdateGroup(ctx context.Context, groupID int32, familyID string, name string) (db.Group, error) {
	group, err := s.groupStore.FindGroupByID(ctx, groupID)
	if err != nil {
		return db.Group{}, err
	}
	if group.FamilyID != familyID {
		return db.Group{}, apperr.NewForbiddenError("error.forbidden", nil)
	}
	if group.IsAdmin {
		return db.Group{}, apperr.NewBadRequestError("error.group.admin_not_editable", nil)
	}
	return s.groupStore.UpdateGroup(ctx, db.UpdateGroupParams{
		Name: name,
		ID:   groupID,
	})
}

func (s *GroupService) DeleteGroup(ctx context.Context, groupID int32, familyID string) error {
	group, err := s.groupStore.FindGroupByID(ctx, groupID)
	if err != nil {
		return err
	}
	if group.FamilyID != familyID {
		return apperr.NewForbiddenError("error.forbidden", nil)
	}
	if group.IsAdmin {
		return apperr.NewBadRequestError("error.group.admin_not_deletable", nil)
	}
	return s.groupStore.DeleteGroup(ctx, groupID)
}
