package services

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/apperr"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store"
)

type GroupService interface {
	GetGroups(ctx context.Context, familyID int32) ([]db.Group, error)
	CreateGroup(ctx context.Context, familyID int32, name string) (db.Group, error)
	UpdateGroup(ctx context.Context, groupID int32, familyID int32, name string) (db.Group, error)
	DeleteGroup(ctx context.Context, groupID int32, familyID int32) error
}

type groupService struct {
	familyStore store.FamilyStore
	groupStore  store.GroupStore
}

func NewGroupService(familyStore store.FamilyStore, groupStore store.GroupStore) GroupService {
	return &groupService{familyStore: familyStore, groupStore: groupStore}
}

func (s *groupService) GetGroups(ctx context.Context, familyID int32) ([]db.Group, error) {
	return s.groupStore.FindGroupsByFamilyID(ctx, familyID)
}

func (s *groupService) CreateGroup(ctx context.Context, familyID int32, name string) (db.Group, error) {
	return s.groupStore.CreateGroup(ctx, db.CreateGroupParams{
		FamilyID: familyID,
		IsAdmin:  false,
		Name:     name,
	})
}

func (s *groupService) UpdateGroup(ctx context.Context, groupID int32, familyID int32, name string) (db.Group, error) {
	group, err := s.groupStore.FindGroupByID(ctx, groupID)
	if err != nil {
		return db.Group{}, err
	}
	if group.FamilyID != familyID {
		return db.Group{}, apperr.NewForbiddenError("error.forbidden", nil)
	}
	return s.groupStore.UpdateGroup(ctx, db.UpdateGroupParams{
		Name: name,
		ID:   groupID,
	})
}

func (s *groupService) DeleteGroup(ctx context.Context, groupID int32, familyID int32) error {
	group, err := s.groupStore.FindGroupByID(ctx, groupID)
	if err != nil {
		return err
	}
	if group.FamilyID != familyID {
		return apperr.NewForbiddenError("error.forbidden", nil)
	}
	return s.groupStore.DeleteGroup(ctx, groupID)
}
