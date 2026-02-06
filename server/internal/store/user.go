package store

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

// UserStore 사용자 데이터 접근 인터페이스
type UserStore interface {
	CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error)
	FindUserByUsername(ctx context.Context, username string) (db.User, error)
}

type userStore struct {
	queries *db.Queries
}

// NewUserStore UserStore 구현체 생성
func NewUserStore(queries *db.Queries) UserStore {
	return &userStore{queries: queries}
}

func (s *userStore) CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	return s.queries.CreateUser(ctx, params)
}

func (s *userStore) FindUserByUsername(ctx context.Context, username string) (db.User, error) {
	return s.queries.FindUserByUsername(ctx, username)
}
