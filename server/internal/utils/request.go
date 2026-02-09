package utils

import (
	"github.com/labstack/echo/v4"
)

// echo context에서 request_id를 가져온다.
func GetRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}
