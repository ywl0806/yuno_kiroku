package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims interface {
	jwt.Claims
}

type AccessTokenClaims struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	GroupId     string `json:"group_id"`
	ClanGroupId string `json:"clan_group_id"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a JWT for the given claims type (AccessToken or RefreshToken).
// expireSeconds is the token TTL in seconds (e.g. 3600 for 1 hour).
func GenerateJWT(claims JWTClaims, secretKey string, expireSeconds int) (string, error) {
	expireAt := time.Now().Add(time.Duration(expireSeconds) * time.Second)

	// Set common registered claims
	switch c := claims.(type) {
	case *AccessTokenClaims:
		c.RegisteredClaims = jwt.RegisteredClaims{
			Issuer:    "app",
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		}
	case *RefreshTokenClaims:
		c.RegisteredClaims = jwt.RegisteredClaims{
			Issuer:    "app",
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// ParseJWT parses a JWT and returns the claims based on the token type.
func ParseJWT(tokenString string, secretKey string, claims JWTClaims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return jwt.ErrInvalidType
	}

	return nil
}
