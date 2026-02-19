package store

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

type InviteTokenStore interface {
	CreateInviteToken(ctx context.Context, params db.CreateInviteTokenParams) (db.InviteToken, error)
	GetInviteTokenByToken(ctx context.Context, token string) (db.InviteToken, error)
	MarkInviteTokenUsed(ctx context.Context, token string) error
}

type inviteTokenStore struct {
	queries *db.Queries
}

func NewInviteTokenStore(queries *db.Queries) InviteTokenStore {
	return &inviteTokenStore{queries: queries}
}

func (s *inviteTokenStore) CreateInviteToken(ctx context.Context, params db.CreateInviteTokenParams) (db.InviteToken, error) {
	return wrapErr(s.queries.CreateInviteToken(ctx, params))
}

func (s *inviteTokenStore) GetInviteTokenByToken(ctx context.Context, token string) (db.InviteToken, error) {
	return wrapErr(s.queries.GetInviteTokenByToken(ctx, token))
}

func (s *inviteTokenStore) MarkInviteTokenUsed(ctx context.Context, token string) error {
	err := s.queries.MarkInviteTokenUsed(ctx, token)
	return mapDBError(err)
}
