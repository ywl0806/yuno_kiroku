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
	lineAuthURL    = "https://access.line.me/oauth2/v2.1/authorize"
	lineTokenURL   = "https://api.line.me/oauth2/v2.1/token"
	lineProfileURL = "https://api.line.me/v2/profile"
)

type LineConfig struct {
	ChannelID     string
	ChannelSecret string
	CallbackURL   string
}

func (c *LineConfig) AuthURL(state string) string {
	v := url.Values{}
	v.Set("response_type", "code")
	v.Set("client_id", c.ChannelID)
	v.Set("redirect_uri", c.CallbackURL)
	v.Set("state", state)
	v.Set("scope", "profile openid")
	return lineAuthURL + "?" + v.Encode()
}

type lineTokenResp struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// ExchangeCode LINE 콜백: code로 토큰 조회
func (c *LineConfig) ExchangeCode(code string) (accessToken string, err error) {
	body := url.Values{}
	body.Set("grant_type", "authorization_code")
	body.Set("code", code)
	body.Set("redirect_uri", c.CallbackURL)
	body.Set("client_id", c.ChannelID)
	body.Set("client_secret", c.ChannelSecret)

	req, err := http.NewRequest(http.MethodPost, lineTokenURL, strings.NewReader(body.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("line token exchange failed: %s %s", resp.Status, string(b))
	}

	var tok lineTokenResp
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		return "", err
	}
	return tok.AccessToken, nil
}

type lineProfile struct {
	UserID        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureURL    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
}

func (c *LineConfig) UserProfile(accessToken string) (userID, displayName string, err error) {
	req, err := http.NewRequest(http.MethodGet, lineProfileURL, nil)
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
		return "", "", fmt.Errorf("line profile failed: %s %s", resp.Status, string(b))
	}

	var p lineProfile
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return "", "", err
	}
	return p.UserID, p.DisplayName, nil
}
