package models

// LoginUserResponse 로그인 성공 시 사용자 정보
type LoginUserResponse struct {
	ID        int32  `json:"id"`
	Username  string `json:"username"`
	FamilyID  int32  `json:"family_id"`
	GroupID   int32  `json:"group_id"`
}

// RedirectURLResponse OAuth 로그인용 리다이렉트 URL 응답
type RedirectURLResponse struct {
	RedirectURL string `json:"redirect_url"`
}
