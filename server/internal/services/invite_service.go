package services

import (
	"context"
	"crypto/rand"
	"database/sql"
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
	familyStore      store.FamilyStore
	groupStore       store.GroupStore
}

func NewInviteService(inviteTokenStore store.InviteTokenStore, familyStore store.FamilyStore, groupStore store.GroupStore) *InviteService {
	return &InviteService{
		inviteTokenStore: inviteTokenStore,
		familyStore:      familyStore,
		groupStore:       groupStore,
	}
}

// CreateInviteToken 초대 토큰 생성. family_id, group_id에 가입할 수 있는 링크용 토큰을 발급한다.
func (s *InviteService) CreateInviteToken(ctx context.Context, familyID, groupID, createdByUserID int32, familyTitle, customFamilyTitle string) (db.InviteToken, error) {
	if err := s.validateFamilyExists(ctx, familyID); err != nil {
		return db.InviteToken{}, err
	}
	if err := s.validateGroupExists(ctx, groupID); err != nil {
		return db.InviteToken{}, err
	}
	token := generateSecureToken(32)
	expiresAt := time.Now().AddDate(0, 0, defaultInviteExpiryDays)
	invite, err := s.inviteTokenStore.CreateInviteToken(ctx, db.CreateInviteTokenParams{
		Token:             token,
		FamilyID:          familyID,
		GroupID:           groupID,
		CreatedByUserID:   createdByUserID,
		ExpiresAt:         expiresAt,
		FamilyTitle:       sql.NullString{String: familyTitle, Valid: familyTitle != ""},
		CustomFamilyTitle: sql.NullString{String: customFamilyTitle, Valid: customFamilyTitle != ""},
	})
	if err != nil {
		return db.InviteToken{}, err
	}
	return invite, nil
}

// ValidateInviteToken 토큰이 유효하면 초대 정보(가족/그룹) 반환. 로그인 전 초대 링크 유효성 확인용.
func (s *InviteService) ValidateInviteToken(ctx context.Context, token string) (familyName, groupName string, familyID, groupID int32, err error) {
	invite, err := s.inviteTokenStore.GetInviteTokenByToken(ctx, token)
	if err != nil || invite.ID == 0 {
		return "", "", 0, 0, apperr.NewAppErrorWithData(apperr.NotFound, "error.invite_token_invalid", nil)
	}
	family, err := s.familyStore.FindFamilyByID(ctx, invite.FamilyID)
	if err != nil || family.ID == 0 {
		return "", "", 0, 0, apperr.NewAppErrorWithData(apperr.NotFound, "error.not_found", map[string]string{"field": "field.family"})
	}
	group, err := s.groupStore.FindGroupByID(ctx, invite.GroupID)
	if err != nil || group.ID == 0 {
		return "", "", 0, 0, apperr.NewAppErrorWithData(apperr.NotFound, "error.not_found", map[string]string{"field": "field.group"})
	}
	return family.Name, group.Name, invite.FamilyID, invite.GroupID, nil
}

// GetValidInviteToken 유효한 초대 토큰 조회 (OAuth/가입 시 사용할 family_id, group_id, family_title 획득용)
func (s *InviteService) GetValidInviteToken(ctx context.Context, token string) (familyID, groupID int32, familyTitle, customFamilyTitle string, err error) {
	invite, err := s.inviteTokenStore.GetInviteTokenByToken(ctx, token)
	if err != nil || invite.ID == 0 {
		return 0, 0, "", "", apperr.NewAppErrorWithData(apperr.NotFound, "error.invite_token_invalid", nil)
	}
	return invite.FamilyID, invite.GroupID, invite.FamilyTitle.String, invite.CustomFamilyTitle.String, nil
}

// MarkInviteTokenUsed 초대 토큰 사용 처리 (한 번만 사용 가능)
func (s *InviteService) MarkInviteTokenUsed(ctx context.Context, token string) error {
	return s.inviteTokenStore.MarkInviteTokenUsed(ctx, token)
}

// generateSecureToken 안전한 토큰 생성
func generateSecureToken(byteLen int) string {
	b := make([]byte, byteLen)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// validateFamilyExists 가족 존재 확인
func (s *InviteService) validateFamilyExists(ctx context.Context, familyID int32) error {
	family, err := s.familyStore.FindFamilyByID(ctx, familyID)
	if err != nil {
		return err
	}
	if family.ID == 0 {
		return apperr.NewAppErrorWithData(apperr.NotFound, "error.not_found", map[string]string{"field": "field.family"})
	}
	return nil
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
