package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/i18n"
	"golang.org/x/text/language"
)

const localeKey = "locale"

// 로케일 설정 미들웨어
// Accept-Language 헤더를 파싱하여 로케일을 설정
func LocaleMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			localeTags, _, _ := language.ParseAcceptLanguage(c.Request().Header.Get("Accept-Language"))
			locales := []string{}
			for _, localeTag := range localeTags {
				locales = append(locales, localeTag.String())
			}
			c.Set(localeKey, locales)
			return next(c)
		}
	}
}

// 로케일 조회
func GetLocale(c echo.Context) []string {
	v, ok := c.Get(localeKey).([]string)
	if !ok || len(v) == 0 {
		return []string{i18n.DefaultLocale}
	}
	return v
}
