package user

import (
	"context"
	"errors"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

type UserService struct {
	queries *db.Queries
}

func NewUserService(queries *db.Queries) *UserService {
	return &UserService{queries: queries}
}

func (s *UserService) CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error) {

	user, err := s.queries.CreateUser(ctx, params)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}

func (s *UserService) FindUserByUsername(ctx context.Context, username string) (db.User, error) {
	user, err := s.queries.FindUserByUsername(ctx, username)

	if err != nil {
		return db.User{}, err
	}
	if user.ID == 0 {
		return db.User{}, errors.New("user not found")
	}
	return user, nil
}

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
	group, err := s.queries.FindGroupByID(ctx, groupID)
	if err != nil {
		return err
	}
	if group.ID == 0 {
		return errors.New("group not found")
	}
	return nil
}

// clan group 존재 체크
func (s *UserService) validateClanGroupExists(ctx context.Context, clanGroupID int32) error {
	clanGroup, err := s.queries.FindClanGroupByID(ctx, clanGroupID)
	if err != nil {
		return err
	}
	if clanGroup.ID == 0 {
		return errors.New("clan group not found")
	}
	return nil
}

// username 중복 체크
func (s *UserService) validateDuplicateUsername(ctx context.Context, username string) error {
	user, err := s.queries.FindUserByUsername(ctx, username)
	if err != nil {
		return err
	}
	if user.ID != 0 {
		return errors.New("username already exists")
	}
	return nil
}
