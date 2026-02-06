package services

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/apperr"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store"
)

type UserService struct {
	userStore  store.UserStore
	groupStore store.GroupStore
}

func NewUserService(userStore store.UserStore, groupStore store.GroupStore) *UserService {
	return &UserService{userStore: userStore, groupStore: groupStore}
}

// 유저 생성
func (s *UserService) CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error) {

	user, err := s.userStore.CreateUser(ctx, params)
	if err != nil {
		return db.User{}, apperr.ClassifyDBError(err)
	}
	return user, nil
}

// 유저 이름으로 유저 조회
func (s *UserService) FindUserByUsername(ctx context.Context, username string) (db.User, error) {
	user, err := s.userStore.FindUserByUsername(ctx, username)
	if err != nil {
		return db.User{}, apperr.ClassifyDBError(err, "user")
	}
	return user, nil
}

// 유저 생성 파라미터 검증
func (s *UserService) ValidateCreateUserParams(ctx context.Context, params db.CreateUserParams) error {
	// Username 중복체크
	err := s.validateDuplicateUsername(ctx, params.Username)
	if err != nil {
		return err
	}
	// group 존재 체크
	err = s.validateGroupExists(ctx, params.GroupID)
	if err != nil {
		return err
	}
	// clan group 존재 체크
	err = s.validateClanGroupExists(ctx, params.ClanGroupID)
	if err != nil {
		return err
	}
	return nil
}

// group 존재 체크
func (s *UserService) validateGroupExists(ctx context.Context, groupID int32) error {
	group, err := s.groupStore.FindGroupByID(ctx, groupID)
	if err != nil {
		return err
	}
	if group.ID == 0 {
		return apperr.NewNotFoundError("group")
	}
	return nil
}

// clan group 존재 체크
func (s *UserService) validateClanGroupExists(ctx context.Context, clanGroupID int32) error {
	clanGroup, err := s.groupStore.FindClanGroupByID(ctx, clanGroupID)
	if err != nil {
		return err
	}
	if clanGroup.ID == 0 {
		return apperr.NewNotFoundError("clan group")
	}
	return nil
}

// username 중복 체크
func (s *UserService) validateDuplicateUsername(ctx context.Context, username string) error {
	user, err := s.userStore.FindUserByUsername(ctx, username)
	if err != nil {
		return err
	}
	if user.ID != 0 {
		return apperr.NewConflictError("")
	}
	return nil
}
