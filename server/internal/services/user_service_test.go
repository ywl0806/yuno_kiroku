package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ywl0806/yuno_kiroku/internal/apperr"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store/mocks"
)

// FindUserByProvider가 기존 유저를 찾으면 그대로 반환하고, CreateUserOAuth는 호출되지 않아야 한다.
func TestFindOrCreateUserOAuth_ExistingUser(t *testing.T) {
	userStore := mocks.NewMockUserStore(t)
	existing := db.User{ID: 1, Username: "line_abc"}
	userStore.EXPECT().FindUserByProvider(mock.Anything, "line", "abc").Return(existing, nil)

	svc := NewUserService(userStore, mocks.NewMockFamilyStore(t), mocks.NewMockGroupStore(t))

	got, err := svc.FindOrCreateUserOAuth(context.Background(), "line", "abc", "name", 1, 1, "", "")

	assert.NoError(t, err)
	assert.Equal(t, existing, got)
	userStore.AssertNotCalled(t, "CreateUserOAuth", mock.Anything, mock.Anything)
}

// 기존 유저가 없으면(sql.ErrNoRows) provider+providerUserID로 username을 만들어 CreateUserOAuth를 호출해야 한다.
func TestFindOrCreateUserOAuth_CreatesNewUser(t *testing.T) {
	userStore := mocks.NewMockUserStore(t)
	userStore.EXPECT().FindUserByProvider(mock.Anything, "line", "abc").Return(db.User{}, sql.ErrNoRows)
	created := db.User{ID: 2, Username: "line_abc"}
	userStore.EXPECT().CreateUserOAuth(mock.Anything, mock.MatchedBy(func(p db.CreateUserOAuthParams) bool {
		return p.Username == "line_abc" && p.FamilyID == 5 && p.GroupID == 6
	})).Return(created, nil)

	svc := NewUserService(userStore, mocks.NewMockFamilyStore(t), mocks.NewMockGroupStore(t))

	got, err := svc.FindOrCreateUserOAuth(context.Background(), "line", "abc", "표시이름", 5, 6, "", "")

	assert.NoError(t, err)
	assert.Equal(t, created, got)
}

// username이 이미 존재하면 Conflict 에러를 반환하고 family/group 조회는 일어나지 않아야 한다.
func TestValidateCreateUserParams_DuplicateUsername(t *testing.T) {
	userStore := mocks.NewMockUserStore(t)
	userStore.EXPECT().FindUserByUsername(mock.Anything, "dup").Return(db.User{ID: 1}, nil)

	svc := NewUserService(userStore, mocks.NewMockFamilyStore(t), mocks.NewMockGroupStore(t))

	err := svc.ValidateCreateUserParams(context.Background(), db.CreateUserParams{Username: "dup"})

	assert.Error(t, err)
	assert.True(t, apperr.IsAppError(err, apperr.Conflict))
}

// 존재하지 않는 familyID로 가입을 시도하면 NotFound 에러를 반환해야 한다.
func TestValidateCreateUserParams_FamilyNotFound(t *testing.T) {
	userStore := mocks.NewMockUserStore(t)
	userStore.EXPECT().FindUserByUsername(mock.Anything, "new").Return(db.User{}, nil)
	familyStore := mocks.NewMockFamilyStore(t)
	familyStore.EXPECT().FindFamilyByID(mock.Anything, int32(99)).Return(db.FindFamilyByIDRow{}, nil)

	svc := NewUserService(userStore, familyStore, mocks.NewMockGroupStore(t))

	err := svc.ValidateCreateUserParams(context.Background(), db.CreateUserParams{Username: "new", FamilyID: 99})

	assert.Error(t, err)
	assert.True(t, apperr.IsAppError(err, apperr.NotFound))
}

// family는 존재하지만 groupID가 존재하지 않으면 NotFound 에러를 반환해야 한다.
func TestValidateCreateUserParams_GroupNotFound(t *testing.T) {
	userStore := mocks.NewMockUserStore(t)
	userStore.EXPECT().FindUserByUsername(mock.Anything, "new").Return(db.User{}, nil)
	familyStore := mocks.NewMockFamilyStore(t)
	familyStore.EXPECT().FindFamilyByID(mock.Anything, int32(1)).Return(db.FindFamilyByIDRow{ID: 1}, nil)
	groupStore := mocks.NewMockGroupStore(t)
	groupStore.EXPECT().FindGroupByID(mock.Anything, int32(404)).Return(db.FindGroupByIDRow{}, nil)

	svc := NewUserService(userStore, familyStore, groupStore)

	err := svc.ValidateCreateUserParams(context.Background(), db.CreateUserParams{Username: "new", FamilyID: 1, GroupID: 404})

	assert.Error(t, err)
	assert.True(t, apperr.IsAppError(err, apperr.NotFound))
}

// username/family/group이 모두 유효하면 에러 없이 통과해야 한다.
func TestValidateCreateUserParams_Success(t *testing.T) {
	userStore := mocks.NewMockUserStore(t)
	userStore.EXPECT().FindUserByUsername(mock.Anything, "new").Return(db.User{}, nil)
	familyStore := mocks.NewMockFamilyStore(t)
	familyStore.EXPECT().FindFamilyByID(mock.Anything, int32(1)).Return(db.FindFamilyByIDRow{ID: 1}, nil)
	groupStore := mocks.NewMockGroupStore(t)
	groupStore.EXPECT().FindGroupByID(mock.Anything, int32(2)).Return(db.FindGroupByIDRow{ID: 2}, nil)

	svc := NewUserService(userStore, familyStore, groupStore)

	err := svc.ValidateCreateUserParams(context.Background(), db.CreateUserParams{Username: "new", FamilyID: 1, GroupID: 2})

	assert.NoError(t, err)
}

// store가 에러(DB 에러 등)를 반환하면 service는 그대로 전파해야 한다 (삼키지 않음).
func TestFindUserByUsername_PropagatesStoreError(t *testing.T) {
	userStore := mocks.NewMockUserStore(t)
	storeErr := errors.New("db down")
	userStore.EXPECT().FindUserByUsername(mock.Anything, "x").Return(db.User{}, storeErr)

	svc := NewUserService(userStore, mocks.NewMockFamilyStore(t), mocks.NewMockGroupStore(t))

	_, err := svc.FindUserByUsername(context.Background(), "x")

	assert.ErrorIs(t, err, storeErr)
}

// GetUserByID는 userStore.FindUserByID 결과를 그대로 전달하는 thin wrapper여야 한다.
func TestGetUserByID_Success(t *testing.T) {
	userStore := mocks.NewMockUserStore(t)
	expected := db.User{ID: 10, FamilyID: 1}
	userStore.EXPECT().FindUserByID(mock.Anything, int32(10), int32(1)).Return(expected, nil)

	svc := NewUserService(userStore, mocks.NewMockFamilyStore(t), mocks.NewMockGroupStore(t))

	got, err := svc.GetUserByID(context.Background(), 10, 1)

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

// 존재하지 않는 유저 조회 시 store 에러가 그대로 전파되어야 한다.
func TestGetUserByID_NotFound(t *testing.T) {
	userStore := mocks.NewMockUserStore(t)
	userStore.EXPECT().FindUserByID(mock.Anything, int32(999), int32(1)).Return(db.User{}, sql.ErrNoRows)

	svc := NewUserService(userStore, mocks.NewMockFamilyStore(t), mocks.NewMockGroupStore(t))

	_, err := svc.GetUserByID(context.Background(), 999, 1)

	assert.ErrorIs(t, err, sql.ErrNoRows)
}

// GetMembers는 같은 familyID에 속한 멤버 목록을 그대로 반환해야 한다.
func TestGetMembers_Success(t *testing.T) {
	userStore := mocks.NewMockUserStore(t)
	members := []db.User{{ID: 1, FamilyID: 7}, {ID: 2, FamilyID: 7}}
	userStore.EXPECT().FindMembersByFamilyID(mock.Anything, int32(7)).Return(members, nil)

	svc := NewUserService(userStore, mocks.NewMockFamilyStore(t), mocks.NewMockGroupStore(t))

	got, err := svc.GetMembers(context.Background(), 7)

	assert.NoError(t, err)
	assert.Equal(t, members, got)
}

// 멤버가 없는 family를 조회하면 빈 슬라이스를 에러 없이 반환해야 한다.
func TestGetMembers_Empty(t *testing.T) {
	userStore := mocks.NewMockUserStore(t)
	userStore.EXPECT().FindMembersByFamilyID(mock.Anything, int32(7)).Return([]db.User{}, nil)

	svc := NewUserService(userStore, mocks.NewMockFamilyStore(t), mocks.NewMockGroupStore(t))

	got, err := svc.GetMembers(context.Background(), 7)

	assert.NoError(t, err)
	assert.Empty(t, got)
}

// name이 비어있지 않으면 UpdateUserNameParams.Name이 Valid:true로 채워져야 한다.
func TestUpdateMe_SetsName(t *testing.T) {
	userStore := mocks.NewMockUserStore(t)
	updated := db.User{ID: 1, Name: sql.NullString{String: "새이름", Valid: true}}
	userStore.EXPECT().UpdateUserName(mock.Anything, mock.MatchedBy(func(p db.UpdateUserNameParams) bool {
		return p.ID == 1 && p.Name.Valid && p.Name.String == "새이름"
	})).Return(updated, nil)

	svc := NewUserService(userStore, mocks.NewMockFamilyStore(t), mocks.NewMockGroupStore(t))

	got, err := svc.UpdateMe(context.Background(), 1, "새이름")

	assert.NoError(t, err)
	assert.Equal(t, updated, got)
}

// name이 빈 문자열이면 NullString이 Valid:false로 전달되어야 한다 (DB에 NULL로 저장).
func TestUpdateMe_EmptyNameIsNull(t *testing.T) {
	userStore := mocks.NewMockUserStore(t)
	userStore.EXPECT().UpdateUserName(mock.Anything, mock.MatchedBy(func(p db.UpdateUserNameParams) bool {
		return p.ID == 1 && !p.Name.Valid
	})).Return(db.User{ID: 1}, nil)

	svc := NewUserService(userStore, mocks.NewMockFamilyStore(t), mocks.NewMockGroupStore(t))

	_, err := svc.UpdateMe(context.Background(), 1, "")

	assert.NoError(t, err)
}

// UpdateMember는 familyTitle/customFamilyTitle을 포함해 UpdateMemberParams를 올바르게 구성해야 한다.
func TestUpdateMember_Success(t *testing.T) {
	userStore := mocks.NewMockUserStore(t)
	updated := db.User{ID: 5, FamilyID: 1, GroupID: 2}
	userStore.EXPECT().UpdateMember(mock.Anything, mock.MatchedBy(func(p db.UpdateMemberParams) bool {
		return p.ID == 5 && p.FamilyID == 1 && p.GroupID == 2 &&
			p.FamilyTitle.Valid && p.FamilyTitle.String == "엄마" &&
			p.CustomFamilyTitle.Valid && p.CustomFamilyTitle.String == "맘"
	})).Return(updated, nil)

	svc := NewUserService(userStore, mocks.NewMockFamilyStore(t), mocks.NewMockGroupStore(t))

	got, err := svc.UpdateMember(context.Background(), 5, 1, 2, "엄마", "맘")

	assert.NoError(t, err)
	assert.Equal(t, updated, got)
}

// CreateUser 성공 시 store가 반환한 유저를 그대로 반환해야 한다.
func TestCreateUser_Success(t *testing.T) {
	userStore := mocks.NewMockUserStore(t)
	params := db.CreateUserParams{Username: "newuser", FamilyID: 1, GroupID: 2}
	created := db.User{ID: 3, Username: "newuser", FamilyID: 1, GroupID: 2}
	userStore.EXPECT().CreateUser(mock.Anything, params).Return(created, nil)

	svc := NewUserService(userStore, mocks.NewMockFamilyStore(t), mocks.NewMockGroupStore(t))

	got, err := svc.CreateUser(context.Background(), params)

	assert.NoError(t, err)
	assert.Equal(t, created, got)
}

// CreateUser 실패(예: username unique 제약 위반) 시 빈 User와 에러를 반환해야 한다.
func TestCreateUser_StoreError(t *testing.T) {
	userStore := mocks.NewMockUserStore(t)
	params := db.CreateUserParams{Username: "dup"}
	storeErr := errors.New("duplicate key value")
	userStore.EXPECT().CreateUser(mock.Anything, params).Return(db.User{}, storeErr)

	svc := NewUserService(userStore, mocks.NewMockFamilyStore(t), mocks.NewMockGroupStore(t))

	got, err := svc.CreateUser(context.Background(), params)

	assert.ErrorIs(t, err, storeErr)
	assert.Equal(t, db.User{}, got)
}
