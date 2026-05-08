package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	kakaoAuthURL   = "https://kauth.kakao.com/oauth/authorize"
	kakaoTokenURL  = "https://kauth.kakao.com/oauth/token"
	kakaoUserURL   = "https://kapi.kakao.com/v2/user/me"
)

type KakaoConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

func (c *KakaoConfig) AuthURL(state string) string {
	v := url.Values{}
	v.Set("response_type", "code")
	v.Set("client_id", c.ClientID)
	v.Set("redirect_uri", c.RedirectURI)
	v.Set("state", state)
	return kakaoAuthURL + "?" + v.Encode()
}

type kakaoTokenResp struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

func (c *KakaoConfig) ExchangeCode(code string) (accessToken string, err error) {
	body := url.Values{}
	body.Set("grant_type", "authorization_code")
	body.Set("client_id", c.ClientID)
	body.Set("redirect_uri", c.RedirectURI)
	body.Set("code", code)
	if c.ClientSecret != "" {
		body.Set("client_secret", c.ClientSecret)
	}

	req, err := http.NewRequest(http.MethodPost, kakaoTokenURL, strings.NewReader(body.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("kakao token exchange failed: %s %s", resp.Status, string(b))
	}

	var tok kakaoTokenResp
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		return "", err
	}
	return tok.AccessToken, nil
}

type kakaoUser struct {
	ID           json.Number `json:"id"`
	Properties   struct {
		Nickname       string `json:"nickname"`
		ProfileImage   string `json:"profile_image"`
		ThumbnailImage string `json:"thumbnail_image"`
	} `json:"properties"`
	KakaoAccount struct {
		Email         string `json:"email"`
		DisplayName   string `json:"display_name"`
		Profile       struct {
			Nickname string `json:"nickname"`
		} `json:"profile"`
	} `json:"kakao_account"`
}

func (c *KakaoConfig) UserProfile(accessToken string) (userID, displayName string, err error) {
	req, err := http.NewRequest(http.MethodGet, kakaoUserURL, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("kakao user failed: %s %s", resp.Status, string(b))
	}

	var u kakaoUser
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return "", "", err
	}
	displayName = u.Properties.Nickname
	if displayName == "" {
		displayName = u.KakaoAccount.Profile.Nickname
	}
	if displayName == "" {
		displayName = u.KakaoAccount.DisplayName
	}
	return u.ID.String(), displayName, nil
}
