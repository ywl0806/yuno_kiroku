package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"

	"github.com/spf13/cast"
	"github.com/ywl0806/yuno_kiroku/internal/consts"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/oauth"
	"github.com/ywl0806/yuno_kiroku/internal/utils/jwt"
	"github.com/ywl0806/yuno_kiroku/pkg/utils"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrLineNotConfigured  = errors.New("LINE login is not configured")
	ErrKakaoNotConfigured = errors.New("Kakao login is not configured")
)

type AuthService struct {
	userService   *UserService
	inviteService *InviteService
	authSecretKey string
	lineConfig    *oauth.LineConfig
	kakaoConfig   *oauth.KakaoConfig
}

func NewAuthService(
	userService *UserService,
	inviteService *InviteService,
	authSecretKey string,
	lineConfig *oauth.LineConfig,
	kakaoConfig *oauth.KakaoConfig,
) *AuthService {
	return &AuthService{
		userService:   userService,
		inviteService: inviteService,
		authSecretKey: authSecretKey,
		lineConfig:    lineConfig,
		kakaoConfig:   kakaoConfig,
	}
}

// LoginResult 로그인 성공 시 반환 데이터
type LoginResult struct {
	User         db.User
	AccessToken  string
	RefreshToken string
}

// Login 아이디/비밀번호 로그인. 성공 시 사용자 정보와 액세스·리프레시 토큰 반환.
func (s *AuthService) Login(ctx context.Context, username, password string) (*LoginResult, error) {
	user, err := s.userService.FindUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user.ID == 0 || !utils.CheckPassword(password, user.Password) {
		return nil, ErrInvalidCredentials
	}
	accessToken, err := s.issueAccessToken(user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.issueRefreshToken(user)
	if err != nil {
		return nil, err
	}
	return &LoginResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ResolveInviteState state가 유효한 초대 토큰이면 해당 group_id, clan_group_id 반환; 아니면 1,1
func (s *AuthService) ResolveInviteState(ctx context.Context, state string) (groupID, clanGroupID int32) {
	groupID, clanGroupID = 1, 1
	if state == "" {
		return groupID, clanGroupID
	}
	gid, cid, err := s.inviteService.GetValidInviteToken(ctx, state)
	if err != nil {
		return groupID, clanGroupID
	}
	return gid, cid
}

// GetLineAuthURL LINE 로그인 URL 반환. 미설정 시 empty string, false
func (s *AuthService) GetLineAuthURL(state string) (url string, configured bool) {
	if s.lineConfig.ChannelID == "" || s.lineConfig.CallbackURL == "" {
		return "", false
	}
	if state == "" {
		state = randomState()
	}
	return s.lineConfig.AuthURL(state), true
}

// ProcessLineCallback LINE 콜백: code로 토큰·프로필 조회 후 유저 생성/조회, 초대 사용 처리 후 유저 반환
func (s *AuthService) ProcessLineCallback(ctx context.Context, code, state string) (*db.User, error) {
	accessToken, err := s.lineConfig.ExchangeCode(code)
	if err != nil {
		log.Println("LINE token exchange error:", err)
		return nil, err
	}
	userID, displayName, err := s.lineConfig.UserProfile(accessToken)
	if err != nil {
		log.Println("LINE profile error:", err)
		return nil, err
	}
	groupID, clanGroupID := s.ResolveInviteState(ctx, state)
	user, err := s.userService.FindOrCreateUserOAuth(ctx, "line", userID, displayName, groupID, clanGroupID)
	if err != nil {
		log.Println("FindOrCreateUserOAuth LINE error:", err)
		return nil, err
	}
	if state != "" {
		_ = s.inviteService.MarkInviteTokenUsed(ctx, state)
	}
	return &user, nil
}

// GetKakaoAuthURL 카카오 로그인 URL 반환. 미설정 시 empty string, false
func (s *AuthService) GetKakaoAuthURL(state string) (url string, configured bool) {
	if s.kakaoConfig.ClientID == "" || s.kakaoConfig.RedirectURI == "" {
		return "", false
	}
	if state == "" {
		state = randomState()
	}
	return s.kakaoConfig.AuthURL(state), true
}

// ProcessKakaoCallback 카카오 콜백: code로 토큰·프로필 조회 후 유저 생성/조회, 초대 사용 처리 후 유저 반환
func (s *AuthService) ProcessKakaoCallback(ctx context.Context, code, state string) (*db.User, error) {
	accessToken, err := s.kakaoConfig.ExchangeCode(code)
	if err != nil {
		log.Println("Kakao token exchange error:", err)
		return nil, err
	}
	userID, displayName, err := s.kakaoConfig.UserProfile(accessToken)
	if err != nil {
		log.Println("Kakao profile error:", err)
		return nil, err
	}
	groupID, clanGroupID := s.ResolveInviteState(ctx, state)
	user, err := s.userService.FindOrCreateUserOAuth(ctx, "kakao", userID, displayName, groupID, clanGroupID)
	if err != nil {
		log.Println("FindOrCreateUserOAuth Kakao error:", err)
		return nil, err
	}
	if state != "" {
		_ = s.inviteService.MarkInviteTokenUsed(ctx, state)
	}
	return &user, nil
}

// IssueOAuthAccessToken OAuth 로그인 유저용 액세스 토큰 발급
func (s *AuthService) IssueOAuthAccessToken(user *db.User) (string, error) {
	return s.issueAccessToken(*user)
}

func (s *AuthService) issueAccessToken(user db.User) (string, error) {
	claims := &jwt.AccessTokenClaims{
		ID:          cast.ToString(user.ID),
		Email:       user.Username,
		GroupId:     cast.ToString(user.GroupID),
		ClanGroupId: cast.ToString(user.ClanGroupID),
	}
	return jwt.GenerateJWT(claims, s.authSecretKey, consts.AccessTokenCookieMaxAge)
}

func (s *AuthService) issueRefreshToken(user db.User) (string, error) {
	claims := &jwt.RefreshTokenClaims{ID: cast.ToString(user.ID)}
	return jwt.GenerateJWT(claims, s.authSecretKey, consts.RefreshTokenCookieMaxAge)
}

func randomState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
