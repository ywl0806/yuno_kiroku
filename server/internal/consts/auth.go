package consts

const (
	RefreshTokenCookieName = "refresh_token"
	AccessTokenCookieName  = "access_token"

	RefreshTokenCookiePath = "/"
	AccessTokenCookiePath  = "/"

	// 30 days in seconds
	RefreshTokenCookieMaxAge = 60 * 60 * 24 * 30
	// 1 hour in seconds
	AccessTokenCookieMaxAge = 60 * 60
)
