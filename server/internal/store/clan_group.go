package store

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

type ClanGroupStore interface {
	FindClanGroupByID(ctx context.Context, id int32) (db.FindClanGroupByIDRow, error)
	CreateClanGroup(ctx context.Context, groupID int32, isAdmin bool, name string) (db.ClanGroup, error)
}

type clanGroupStore struct {
	queries *db.Queries
}

func NewClanGroupStore(queries *db.Queries) ClanGroupStore {
	return &clanGroupStore{queries: queries}
}

func (s *clanGroupStore) FindClanGroupByID(ctx context.Context, id int32) (db.FindClanGroupByIDRow, error) {
	return wrapErr(s.queries.FindClanGroupByID(ctx, id))
}

func (s *clanGroupStore) CreateClanGroup(ctx context.Context, groupID int32, isAdmin bool, name string) (db.ClanGroup, error) {
	return wrapErr(s.queries.CreateClanGroup(ctx, db.CreateClanGroupParams{
		GroupID: groupID,
		IsAdmin: isAdmin,
		Name:    name,
	}))
}
