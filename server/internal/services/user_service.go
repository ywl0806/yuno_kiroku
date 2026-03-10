package services

import (
	"context"
	"database/sql"

	"github.com/ywl0806/yuno_kiroku/internal/apperr"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store"
)

type UserService struct {
	userStore   store.UserStore
	familyStore store.FamilyStore
	groupStore  store.GroupStore
}

func NewUserService(userStore store.UserStore, familyStore store.FamilyStore, groupStore store.GroupStore) *UserService {
	return &UserService{userStore: userStore, familyStore: familyStore, groupStore: groupStore}
}

func (s *UserService) GetMembers(ctx context.Context, familyID int32) ([]db.FindMembersByFamilyIDRow, error) {
	return s.userStore.FindMembersByFamilyID(ctx, familyID)
}

// 유저 생성
func (s *UserService) CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error) {

	user, err := s.userStore.CreateUser(ctx, params)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}

// 유저 이름으로 유저 조회
func (s *UserService) FindUserByUsername(ctx context.Context, username string) (db.User, error) {
	user, err := s.userStore.FindUserByUsername(ctx, username)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}

// FindUserByProvider 소셜 로그인 provider로 유저 조회
func (s *UserService) FindUserByProvider(ctx context.Context, provider, providerUserID string) (db.User, error) {
	return s.userStore.FindUserByProvider(ctx, provider, providerUserID)
}

// FindOrCreateUserOAuth 소셜 로그인 유저 조회 또는 생성. familyID/groupID로 가입할 가족·그룹을 지정한다 (미지정 시 1,1).
func (s *UserService) FindOrCreateUserOAuth(ctx context.Context, provider, providerUserID, displayName string, familyID, groupID int32) (db.User, error) {
	user, err := s.userStore.FindUserByProvider(ctx, provider, providerUserID)
	if err == nil && user.ID != 0 {
		return user, nil
	}
	if familyID == 0 {
		familyID = 1
	}
	if groupID == 0 {
		groupID = 1
	}
	username := provider + "_" + providerUserID
	params := db.CreateUserOAuthParams{
		Name:           sql.NullString{String: displayName, Valid: displayName != ""},
		Username:       username,
		Password:       "",
		FamilyID:       familyID,
		GroupID:        groupID,
		Provider:       sql.NullString{String: provider, Valid: true},
		ProviderUserID: sql.NullString{String: providerUserID, Valid: true},
	}
	return s.userStore.CreateUserOAuth(ctx, params)
}

// 유저 생성 파라미터 검증
func (s *UserService) ValidateCreateUserParams(ctx context.Context, params db.CreateUserParams) error {
	err := s.validateDuplicateUsername(ctx, params.Username)
	if err != nil {
		return err
	}
	err = s.validateFamilyExists(ctx, params.FamilyID)
	if err != nil {
		return err
	}
	err = s.validateGroupExists(ctx, params.GroupID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) validateFamilyExists(ctx context.Context, familyID int32) error {
	family, err := s.familyStore.FindFamilyByID(ctx, familyID)
	if err != nil {
		return err
	}
	if family.ID == 0 {
		return apperr.NewAppErrorWithData(apperr.NotFound, "error.not_found", map[string]string{"field": "field.family"})
	}
	return nil
}

func (s *UserService) validateGroupExists(ctx context.Context, groupID int32) error {
	group, err := s.groupStore.FindGroupByID(ctx, groupID)
	if err != nil {
		return err
	}
	if group.ID == 0 {
		return apperr.NewAppErrorWithData(apperr.NotFound, "error.not_found", map[string]string{"field": "field.group"})
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
		return apperr.NewAppErrorWithData(apperr.Conflict, "error.conflict.username", nil)
	}
	return nil
}
