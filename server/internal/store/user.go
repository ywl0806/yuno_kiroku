package store

import (
	"context"
	"database/sql"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

// UserStore 사용자 데이터 접근 인터페이스
type UserStore interface {
	CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error)
	CreateUserOAuth(ctx context.Context, params db.CreateUserOAuthParams) (db.User, error)
	FindUserByUsername(ctx context.Context, username string) (db.User, error)
	FindUserByProvider(ctx context.Context, provider, providerUserID string) (db.User, error)
	FindMembersByFamilyID(ctx context.Context, familyID string) ([]db.User, error)
	FindUserByID(ctx context.Context, id string) (db.User, error)
	FindUserByIDAndFamilyID(ctx context.Context, id string, familyID string) (db.User, error)
	UpdateUserName(ctx context.Context, arg db.UpdateUserNameParams) (db.User, error)
	UpdateMember(ctx context.Context, arg db.UpdateMemberParams) (db.User, error)
}

type userStore struct {
	queries *db.Queries
}

// NewUserStore UserStore 구현체 생성
func NewUserStore(queries *db.Queries) UserStore {
	return &userStore{queries: queries}
}

func (s *userStore) CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	return wrapErr(s.queries.CreateUser(ctx, params))
}

func (s *userStore) CreateUserOAuth(ctx context.Context, params db.CreateUserOAuthParams) (db.User, error) {
	return wrapErr(s.queries.CreateUserOAuth(ctx, params))
}

func (s *userStore) FindUserByUsername(ctx context.Context, username string) (db.User, error) {
	return wrapErr(s.queries.FindUserByUsername(ctx, username))
}

func (s *userStore) FindMembersByFamilyID(ctx context.Context, familyID string) ([]db.User, error) {
	return wrapErr(s.queries.FindMembersByFamilyID(ctx, familyID))
}

func (s *userStore) FindUserByProvider(ctx context.Context, provider, providerUserID string) (db.User, error) {
	return wrapErr(s.queries.FindUserByProvider(ctx, db.FindUserByProviderParams{
		Provider:       sql.NullString{String: provider, Valid: true},
		ProviderUserID: sql.NullString{String: providerUserID, Valid: true},
	}))
}

func (s *userStore) FindUserByID(ctx context.Context, id string) (db.User, error) {
	return wrapErr(s.queries.FindUserByID(ctx, id))
}

func (s *userStore) FindUserByIDAndFamilyID(ctx context.Context, id string, familyID string) (db.User, error) {
	return wrapErr(s.queries.FindUserByIDAndFamilyID(ctx, db.FindUserByIDAndFamilyIDParams{
		ID:       id,
		FamilyID: familyID,
	}))
}

func (s *userStore) UpdateUserName(ctx context.Context, arg db.UpdateUserNameParams) (db.User, error) {
	return wrapErr(s.queries.UpdateUserName(ctx, arg))
}

func (s *userStore) UpdateMember(ctx context.Context, arg db.UpdateMemberParams) (db.User, error) {
	return wrapErr(s.queries.UpdateMember(ctx, arg))
}
