package services

import (
	"context"
	"database/sql"

	"github.com/ywl0806/yuno_kiroku/internal/apperr"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store"
)

type UserService interface {
	GetUserByID(ctx context.Context, userID, familyID int32) (db.User, error)
	GetMembers(ctx context.Context, familyID int32) ([]db.User, error)
	UpdateMe(ctx context.Context, userID int32, name string) (db.User, error)
	UpdateMember(ctx context.Context, memberID, familyID, groupID int32, familyTitle, customFamilyTitle string) (db.User, error)
	CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error)
	FindUserByUsername(ctx context.Context, username string) (db.User, error)
	FindUserByProvider(ctx context.Context, provider, providerUserID string) (db.User, error)
	FindOrCreateUserOAuth(ctx context.Context, provider, providerUserID, displayName string, familyID, groupID int32, familyTitle, customFamilyTitle string) (db.User, error)
	ValidateCreateUserParams(ctx context.Context, params db.CreateUserParams) error
}

type userService struct {
	userStore   store.UserStore
	familyStore store.FamilyStore
	groupStore  store.GroupStore
}

func NewUserService(userStore store.UserStore, familyStore store.FamilyStore, groupStore store.GroupStore) UserService {
	return &userService{userStore: userStore, familyStore: familyStore, groupStore: groupStore}
}

func (s *userService) GetUserByID(ctx context.Context, userID, familyID int32) (db.User, error) {
	return s.userStore.FindUserByID(ctx, userID, familyID)
}

func (s *userService) GetMembers(ctx context.Context, familyID int32) ([]db.User, error) {
	return s.userStore.FindMembersByFamilyID(ctx, familyID)
}

func (s *userService) UpdateMe(ctx context.Context, userID int32, name string) (db.User, error) {
	return s.userStore.UpdateUserName(ctx, db.UpdateUserNameParams{
		Name: sql.NullString{String: name, Valid: name != ""},
		ID:   userID,
	})
}

func (s *userService) UpdateMember(ctx context.Context, memberID, familyID, groupID int32, familyTitle, customFamilyTitle string) (db.User, error) {
	return s.userStore.UpdateMember(ctx, db.UpdateMemberParams{
		GroupID:           groupID,
		FamilyTitle:       sql.NullString{String: familyTitle, Valid: familyTitle != ""},
		CustomFamilyTitle: sql.NullString{String: customFamilyTitle, Valid: customFamilyTitle != ""},
		ID:                memberID,
		FamilyID:          familyID,
	})
}

// 유저 생성
func (s *userService) CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error) {

	user, err := s.userStore.CreateUser(ctx, params)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}

// 유저 이름으로 유저 조회
func (s *userService) FindUserByUsername(ctx context.Context, username string) (db.User, error) {
	user, err := s.userStore.FindUserByUsername(ctx, username)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}

// FindUserByProvider 소셜 로그인 provider로 유저 조회
func (s *userService) FindUserByProvider(ctx context.Context, provider, providerUserID string) (db.User, error) {
	return s.userStore.FindUserByProvider(ctx, provider, providerUserID)
}

// FindOrCreateUserOAuth 소셜 로그인 유저 조회 또는 생성. familyID/groupID로 가입할 가족·그룹을 지정한다
func (s *userService) FindOrCreateUserOAuth(ctx context.Context, provider, providerUserID, displayName string, familyID, groupID int32, familyTitle, customFamilyTitle string) (db.User, error) {
	user, err := s.userStore.FindUserByProvider(ctx, provider, providerUserID)
	if err == nil && user.ID != 0 {
		return user, nil
	}

	username := provider + "_" + providerUserID
	params := db.CreateUserOAuthParams{
		Name:              sql.NullString{String: displayName, Valid: displayName != ""},
		Username:          username,
		Password:          "",
		FamilyID:          familyID,
		GroupID:           groupID,
		Provider:          sql.NullString{String: provider, Valid: true},
		ProviderUserID:    sql.NullString{String: providerUserID, Valid: true},
		FamilyTitle:       sql.NullString{String: familyTitle, Valid: familyTitle != ""},
		CustomFamilyTitle: sql.NullString{String: customFamilyTitle, Valid: customFamilyTitle != ""},
	}
	return s.userStore.CreateUserOAuth(ctx, params)
}

// 유저 생성 파라미터 검증
func (s *userService) ValidateCreateUserParams(ctx context.Context, params db.CreateUserParams) error {
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

func (s *userService) validateFamilyExists(ctx context.Context, familyID int32) error {
	family, err := s.familyStore.FindFamilyByID(ctx, familyID)
	if err != nil {
		return err
	}
	if family.ID == 0 {
		return apperr.NewAppErrorWithData(apperr.NotFound, "error.not_found", map[string]string{"field": "field.family"})
	}
	return nil
}

func (s *userService) validateGroupExists(ctx context.Context, groupID int32) error {
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
func (s *userService) validateDuplicateUsername(ctx context.Context, username string) error {
	user, err := s.userStore.FindUserByUsername(ctx, username)
	if err != nil {
		return err
	}
	if user.ID != 0 {
		return apperr.NewAppErrorWithData(apperr.Conflict, "error.conflict.username", nil)
	}
	return nil
}
