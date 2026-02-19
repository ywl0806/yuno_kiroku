package models

import "time"

// CreateInviteResponse 초대 토큰 생성 응답
type CreateInviteResponse struct {
	Token     string    `json:"token"`
	InviteURL string    `json:"invite_url"`
	ExpiresAt time.Time `json:"expires_at"`
}

// ValidateInviteResponse 초대 토큰 검증 응답
type ValidateInviteResponse struct {
	Valid         bool   `json:"valid"`
	GroupName     string `json:"group_name"`
	ClanGroupName string `json:"clan_group_name"`
	GroupID       int32  `json:"group_id"`
	ClanGroupID   int32  `json:"clan_group_id"`
}
