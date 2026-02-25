package store

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

// GroupStore 그룹/클랜그룹 데이터 접근 인터페이스
type GroupStore interface {
	FindGroupByID(ctx context.Context, id int32) (db.FindGroupByIDRow, error)
	FindClanGroupByID(ctx context.Context, id int32) (db.FindClanGroupByIDRow, error)
	CreateGroup(ctx context.Context, name string) (db.Group, error)
}

type groupStore struct {
	queries *db.Queries
}

// NewGroupStore GroupStore 구현체 생성
func NewGroupStore(queries *db.Queries) GroupStore {
	return &groupStore{queries: queries}
}

func (s *groupStore) FindGroupByID(ctx context.Context, id int32) (db.FindGroupByIDRow, error) {
	row, err := s.queries.FindGroupByID(ctx, id)
	return wrapErr(row, err, "field.group")
}

func (s *groupStore) FindClanGroupByID(ctx context.Context, id int32) (db.FindClanGroupByIDRow, error) {
	row, err := s.queries.FindClanGroupByID(ctx, id)
	return wrapErr(row, err, "field.clan_group")
}

func (s *groupStore) CreateGroup(ctx context.Context, name string) (db.Group, error) {
	return wrapErr(s.queries.CreateGroup(ctx, name))
}
