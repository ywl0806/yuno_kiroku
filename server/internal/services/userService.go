package services

import (
	"context"
	"database/sql"

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

// FindOrCreateUserOAuth 소셜 로그인 유저 조회 또는 생성. groupID/clanGroupID로 가입할 그룹을 지정한다 (미지정 시 1,1).
func (s *UserService) FindOrCreateUserOAuth(ctx context.Context, provider, providerUserID, displayName string, groupID, clanGroupID int32) (db.User, error) {
	user, err := s.userStore.FindUserByProvider(ctx, provider, providerUserID)
	if err == nil && user.ID != 0 {
		return user, nil
	}
	if groupID == 0 {
		groupID = 1
	}
	if clanGroupID == 0 {
		clanGroupID = 1
	}
	username := provider + "_" + providerUserID
	params := db.CreateUserOAuthParams{
		Name:           sql.NullString{String: displayName, Valid: displayName != ""},
		Username:       username,
		Password:       "",
		GroupID:        groupID,
		ClanGroupID:    clanGroupID,
		Provider:       sql.NullString{String: provider, Valid: true},
		ProviderUserID: sql.NullString{String: providerUserID, Valid: true},
	}
	return s.userStore.CreateUserOAuth(ctx, params)
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
		return apperr.NewAppErrorWithData(apperr.NotFound, "error.not_found", map[string]string{"field": "field.group"})
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
		return apperr.NewAppErrorWithData(apperr.NotFound, "error.not_found", map[string]string{"field": "field.clan_group"})
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
