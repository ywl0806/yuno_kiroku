package services

import (
	"context"
	"database/sql"

	"github.com/ywl0806/yuno_kiroku/internal/apperr"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store"
)

type UserService struct {
	userStore                 store.UserStore
	familyStore               store.FamilyStore
	groupStore                store.GroupStore
	albumGroupPermissionStore store.AlbumGroupPermissionStore
}

func NewUserService(userStore store.UserStore, familyStore store.FamilyStore, groupStore store.GroupStore, albumGroupPermissionStore store.AlbumGroupPermissionStore) *UserService {
	return &UserService{userStore: userStore, familyStore: familyStore, groupStore: groupStore, albumGroupPermissionStore: albumGroupPermissionStore}
}

func (s *UserService) GetUserByID(ctx context.Context, userID string) (db.User, error) {
	return s.userStore.FindUserByID(ctx, userID)
}

func (s *UserService) GetUserByIDAndFamilyID(ctx context.Context, userID, familyID string) (db.User, error) {
	return s.userStore.FindUserByIDAndFamilyID(ctx, userID, familyID)
}

func (s *UserService) GetMembers(ctx context.Context, familyID string) ([]db.User, error) {
	return s.userStore.FindMembersByFamilyID(ctx, familyID)
}

func (s *UserService) UpdateMe(ctx context.Context, userID string, name string) (db.User, error) {
	return s.userStore.UpdateUserName(ctx, db.UpdateUserNameParams{
		Name: sql.NullString{String: name, Valid: name != ""},
		ID:   userID,
	})
}

func (s *UserService) UpdateMember(ctx context.Context, memberID string, familyID string, groupID int32, familyTitle, customFamilyTitle string) (db.User, error) {
	return s.userStore.UpdateMember(ctx, db.UpdateMemberParams{
		GroupID:           groupID,
		FamilyTitle:       sql.NullString{String: familyTitle, Valid: familyTitle != ""},
		CustomFamilyTitle: sql.NullString{String: customFamilyTitle, Valid: customFamilyTitle != ""},
		ID:                memberID,
		FamilyID:          familyID,
	})
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

// FindOrCreateUserOAuth 소셜 로그인 유저 조회 또는 생성. familyID/groupID로 가입할 가족·그룹을 지정한다
func (s *UserService) FindOrCreateUserOAuth(ctx context.Context, provider, providerUserID, displayName string, familyID string, groupID int32, familyTitle, customFamilyTitle string) (db.User, error) {
	user, err := s.userStore.FindUserByProvider(ctx, provider, providerUserID)
	if err == nil && user.ID != "" {
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

func (s *UserService) validateFamilyExists(ctx context.Context, familyID string) error {
	_, err := s.familyStore.FindFamilyByID(ctx, familyID)
	return err
}

func (s *UserService) GetWritableAlbumIDs(ctx context.Context, groupID int32) ([]string, error) {
	ids, err := s.albumGroupPermissionStore.GetWritableAlbumIDsByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if ids == nil {
		return []string{}, nil
	}
	return ids, nil
}

func (s *UserService) GetGroupIsAdmin(ctx context.Context, groupID int32) (bool, error) {
	group, err := s.groupStore.FindGroupByID(ctx, groupID)
	if err != nil {
		return false, err
	}
	return group.IsAdmin, nil
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
	if user.ID != "" {
		return apperr.NewAppErrorWithData(apperr.Conflict, "error.conflict.username", nil)
	}
	return nil
}
