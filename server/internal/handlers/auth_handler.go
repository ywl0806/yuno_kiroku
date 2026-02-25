package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/internal/apperr"
	"github.com/ywl0806/yuno_kiroku/internal/consts"
	"github.com/ywl0806/yuno_kiroku/internal/handlers/models"
	"github.com/ywl0806/yuno_kiroku/internal/services"

	_ "github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService *services.AuthService
	frontURL    string
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	frontURL := viper.GetString("FRONT_URL")
	if frontURL == "" {
		frontURL = "http://localhost:5173"
	}
	return &AuthHandler{
		authService: authService,
		frontURL:    frontURL,
	}
}

// LoginRequest 로그인 요청
type LoginRequest struct {
	Username string `json:"username" validate:"required" example:"admin"`
	Password string `json:"password" validate:"required" example:"password"`
}

// LoginResponse 로그인 성공 응답
type LoginResponse struct {
	User  models.LoginUserResponse `json:"user"`
	Token string                   `json:"token"`
}

// @Description 아이디/비밀번호 로그인. 성공 시 액세스 토큰과 리프레시 토큰(쿠키) 반환.
//
// @Summary User Login
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 401 "Invalid username or password"
// @Failure 500 "Failed to generate token"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	result, err := h.authService.Login(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			return apperr.NewUnauthorizedError("message.invalid_credentials", nil)
		}
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:     consts.RefreshTokenCookieName,
		Value:    result.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		MaxAge:   consts.RefreshTokenCookieMaxAge,
	})

	return c.JSON(http.StatusOK, LoginResponse{
		User: models.LoginUserResponse{
			ID:          result.User.ID,
			Username:    result.User.Username,
			GroupID:     result.User.GroupID,
			ClanGroupID: result.User.ClanGroupID,
		},
		Token: result.AccessToken,
	})
}

// @Description LINE 로그인 페이지로 리다이렉트. query invite_token이 있으면 state로 넘겨 콜백에서 초대 그룹 적용.
//
// @Summary LINE Login Redirect
// @Tags Auth
// @Accept json
// @Produce json
// @Param invite_token query string true "Invite Token (required)"
// @Success 302 "Redirect to LINE authorization page"
// @Failure 503 "LINE login is not configured"
// @Router /auth/line [get]
func (h *AuthHandler) LineLoginRedirect(c echo.Context) error {
	state := c.QueryParam("invite_token")
	if state != "" {
		return apperr.NewBadRequestError("message.domain.login-invited-user-only", nil)
	}
	url, configured := h.authService.GetLineAuthURL(state)
	if !configured {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "LINE login is not configured")
	}
	return c.Redirect(http.StatusFound, url)
}

// @Description LINE 로그인 콜백. code로 토큰·프로필 조회 후 유저 생성/조회 및 JWT 발급, 프론트 로그인 콜백 URL로 리다이렉트.
//
// @Summary LINE Login Callback
// @Tags Auth
// @Accept json
// @Produce json
// @Param code query string true "Authorization code from LINE"
// @Param state query string true "State (invite_token required)"
// @Success 302 "Redirect to front login/callback with token"
// @Failure 302 "Redirect to front /login with error query"
// @Router /auth/line/callback [get]
func (h *AuthHandler) LineCallback(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return c.Redirect(http.StatusFound, h.frontURL+"/login?error=missing_code")
	}
	state := c.QueryParam("state")
	if state != "" {
		return apperr.NewBadRequestError("message.domain.login-invited-user-only", nil)
	}
	user, err := h.authService.ProcessLineCallback(c.Request().Context(), code, state)
	if err != nil {
		return c.Redirect(http.StatusFound, h.frontURL+"/login?error=line_token")
	}

	token, err := h.authService.IssueOAuthAccessToken(user)
	if err != nil {
		return c.Redirect(http.StatusFound, h.frontURL+"/login?error=token")
	}
	return c.Redirect(http.StatusFound, h.frontURL+"/login/callback?token="+token)
}

// @Description 카카오 로그인 페이지로 리다이렉트. query invite_token이 있으면 state로 넘겨 콜백에서 초대 그룹 적용.
//
// @Summary Kakao Login Redirect
// @Tags Auth
// @Accept json
// @Produce json
// @Param invite_token query string true "Invite Token (required)"
// @Success 302 "Redirect to Kakao authorization page"
// @Failure 503 "Kakao login is not configured"
// @Router /auth/kakao [get]
func (h *AuthHandler) KakaoLoginRedirect(c echo.Context) error {
	state := c.QueryParam("invite_token")

	if state != "" {
		return apperr.NewBadRequestError("message.domain.login-invited-user-only", nil)
	}
	url, configured := h.authService.GetKakaoAuthURL(state)
	if !configured {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Kakao login is not configured")
	}
	return c.Redirect(http.StatusFound, url)
}

// @Description 카카오 로그인 콜백. code로 토큰·프로필 조회 후 유저 생성/조회 및 JWT 발급, 프론트 로그인 콜백 URL로 리다이렉트.
//
// @Summary Kakao Login Callback
// @Tags Auth
// @Accept json
// @Produce json
// @Param code query string true "Authorization code from Kakao"
// @Param state query string true "State (invite_token required)"
// @Success 302 "Redirect to front login/callback with token"
// @Failure 302 "Redirect to front /login with error query"
// @Router /auth/kakao/callback [get]
func (h *AuthHandler) KakaoCallback(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return c.Redirect(http.StatusFound, h.frontURL+"/login?error=missing_code")
	}
	state := c.QueryParam("state")
	if state != "" {
		return apperr.NewBadRequestError("message.domain.login-invited-user-only", nil)
	}
	user, err := h.authService.ProcessKakaoCallback(c.Request().Context(), code, state)
	if err != nil {
		return c.Redirect(http.StatusFound, h.frontURL+"/login?error=kakao_token")
	}

	token, err := h.authService.IssueOAuthAccessToken(user)
	if err != nil {
		return c.Redirect(http.StatusFound, h.frontURL+"/login?error=token")
	}
	return c.Redirect(http.StatusFound, h.frontURL+"/login/callback?token="+token)
}
