package cloudfront

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	cfsign "github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign"
	"github.com/labstack/echo/v4"
)

const CookieTTL = 12 * time.Hour

type Signer struct {
	keyPairID    string
	cookieSigner *cfsign.CookieSigner
	domain       string // media CloudFront 도메인 (예: media.yuno.jp)
	cookieDomain string // Set-Cookie Domain 속성에 사용할 상위 도메인 (예: yuno.jp)
	secure       bool
}

// parentDomain은 "sub.example.com"에서 "example.com"을 추출한다.
// API 서버(api.yuno.jp)가 미디어 도메인(media.yuno.jp)에 직접 쿠키를 설정할 수 없으므로
// 공통 상위 도메인을 사용해야 브라우저가 쿠키를 수락한다.
func parentDomain(domain string) string {
	if _, after, ok := strings.Cut(domain, "."); ok {
		return after
	}
	return domain
}

func NewSigner(keyPairID, privateKeyPEM, domain string, secure bool) (*Signer, error) {
	privateKeyPEM = strings.ReplaceAll(privateKeyPEM, `\n`, "\n")

	privKey, err := cfsign.LoadPEMPrivKey(strings.NewReader(privateKeyPEM))
	if err != nil {
		// PKCS8 형식도 시도
		signer, err2 := cfsign.LoadPEMPrivKeyPKCS8AsSigner(strings.NewReader(privateKeyPEM))
		if err2 != nil {
			return nil, fmt.Errorf("cloudfront: failed to load private key: %w", err)
		}
		return &Signer{
			keyPairID:    keyPairID,
			cookieSigner: cfsign.NewCookieSigner(keyPairID, signer),
			domain:       domain,
			cookieDomain: parentDomain(domain),
			secure:       secure,
		}, nil
	}

	return &Signer{
		keyPairID:    keyPairID,
		cookieSigner: cfsign.NewCookieSigner(keyPairID, privKey),
		domain:       domain,
		cookieDomain: parentDomain(domain),
		secure:       secure,
	}, nil
}

// SetFamilyCookies family의 media/* 경로에 대한 CloudFront signed cookie를 응답에 설정
func (s *Signer) SetFamilyCookies(c echo.Context, familyID string) error {
	resourceURL := fmt.Sprintf("https://%s/media/%s/*", s.domain, familyID)
	expiry := time.Now().Add(CookieTTL)

	cookies, err := s.cookieSigner.Sign(resourceURL, expiry)
	if err != nil {
		return fmt.Errorf("cloudfront: failed to sign cookies: %w", err)
	}

	cookiePath := fmt.Sprintf("/media/%s/", familyID)
	for _, cookie := range cookies {
		cookie.Path = cookiePath
		cookie.Domain = s.cookieDomain
		cookie.Secure = s.secure
		cookie.HttpOnly = true
		cookie.SameSite = http.SameSiteStrictMode
		c.SetCookie(cookie)
	}
	return nil
}

// ClearFamilyCookies CloudFront signed cookie를 만료 처리
func (s *Signer) ClearFamilyCookies(c echo.Context, familyID string) {
	cookiePath := fmt.Sprintf("/media/%s/", familyID)
	for _, name := range []string{cfsign.CookiePolicyName, cfsign.CookieSignatureName, cfsign.CookieKeyIDName} {
		c.SetCookie(&http.Cookie{
			Name:     name,
			Value:    "",
			Path:     cookiePath,
			Domain:   s.cookieDomain,
			MaxAge:   -1,
			Secure:   s.secure,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})
	}
}
