package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/ywl0806/yuno_kiroku/internal/apperr"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store"
)

// 기본 초대 만료 기간 (7일)
const defaultInviteExpiryDays = 7

type InviteService struct {
	inviteTokenStore store.InviteTokenStore
	groupStore       store.GroupStore
}

func NewInviteService(inviteTokenStore store.InviteTokenStore, groupStore store.GroupStore) *InviteService {
	return &InviteService{
		inviteTokenStore: inviteTokenStore,
		groupStore:       groupStore,
	}
}

// CreateInviteToken 초대 토큰 생성. group_id, clan_group_id에 가입할 수 있는 링크용 토큰을 발급한다.
func (s *InviteService) CreateInviteToken(ctx context.Context, groupID, clanGroupID, createdByUserID int32) (db.InviteToken, error) {
	if err := s.validateGroupExists(ctx, groupID); err != nil {
		return db.InviteToken{}, err
	}
	if err := s.validateClanGroupExists(ctx, clanGroupID); err != nil {
		return db.InviteToken{}, err
	}
	token := generateSecureToken(32)
	expiresAt := time.Now().AddDate(0, 0, defaultInviteExpiryDays)
	invite, err := s.inviteTokenStore.CreateInviteToken(ctx, db.CreateInviteTokenParams{
		Token:           token,
		GroupID:         groupID,
		ClanGroupID:     clanGroupID,
		CreatedByUserID: createdByUserID,
		ExpiresAt:       expiresAt,
	})
	if err != nil {
		return db.InviteToken{}, err
	}
	return invite, nil
}

// ValidateInviteToken 토큰이 유효하면 초대 정보(그룹/클랜) 반환. 로그인 전 초대 링크 유효성 확인용.
func (s *InviteService) ValidateInviteToken(ctx context.Context, token string) (groupName, clanGroupName string, groupID, clanGroupID int32, err error) {
	invite, err := s.inviteTokenStore.GetInviteTokenByToken(ctx, token)
	if err != nil || invite.ID == 0 {
		return "", "", 0, 0, apperr.NewAppErrorWithData(apperr.NotFound, "error.invite_token_invalid", nil)
	}
	group, err := s.groupStore.FindGroupByID(ctx, invite.GroupID)
	if err != nil || group.ID == 0 {
		return "", "", 0, 0, apperr.NewAppErrorWithData(apperr.NotFound, "error.not_found", map[string]string{"field": "field.group"})
	}
	clanGroup, err := s.groupStore.FindClanGroupByID(ctx, invite.ClanGroupID)
	if err != nil || clanGroup.ID == 0 {
		return "", "", 0, 0, apperr.NewAppErrorWithData(apperr.NotFound, "error.not_found", map[string]string{"field": "field.clan_group"})
	}
	return group.Name, clanGroup.Name, invite.GroupID, invite.ClanGroupID, nil
}

// GetValidInviteToken 유효한 초대 토큰 조회 (OAuth/가입 시 사용할 group_id, clan_group_id 획득용)
func (s *InviteService) GetValidInviteToken(ctx context.Context, token string) (groupID, clanGroupID int32, err error) {
	invite, err := s.inviteTokenStore.GetInviteTokenByToken(ctx, token)
	if err != nil || invite.ID == 0 {
		return 0, 0, apperr.NewAppErrorWithData(apperr.NotFound, "error.invite_token_invalid", nil)
	}
	return invite.GroupID, invite.ClanGroupID, nil
}

// MarkInviteTokenUsed 초대 토큰 사용 처리 (한 번만 사용 가능)
func (s *InviteService) MarkInviteTokenUsed(ctx context.Context, token string) error {
	return s.inviteTokenStore.MarkInviteTokenUsed(ctx, token)
}

func generateSecureToken(byteLen int) string {
	b := make([]byte, byteLen)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *InviteService) validateGroupExists(ctx context.Context, groupID int32) error {
	group, err := s.groupStore.FindGroupByID(ctx, groupID)
	if err != nil {
		return err
	}
	if group.ID == 0 {
		return apperr.NewAppErrorWithData(apperr.NotFound, "error.not_found", map[string]string{"field": "field.group"})
	}
	return nil
}

func (s *InviteService) validateClanGroupExists(ctx context.Context, clanGroupID int32) error {
	clanGroup, err := s.groupStore.FindClanGroupByID(ctx, clanGroupID)
	if err != nil {
		return err
	}
	if clanGroup.ID == 0 {
		return apperr.NewAppErrorWithData(apperr.NotFound, "error.not_found", map[string]string{"field": "field.clan_group"})
	}
	return nil
}
