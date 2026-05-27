package middlewares

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/utils"
)

// RequestLogger는 Echo 기본 Logger 미들웨어를 대체한다.
// request_id는 RequestIDWithConfig 미들웨어에서 이미 context에 심겨 있다.
func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()
			ctx := req.Context()

			err := next(c)

			slog.InfoContext(ctx, "http",
				"method", req.Method,
				"path", req.URL.Path,
				"status", c.Response().Status,
				"duration_ms", time.Since(start).Milliseconds(),
				"request_id", utils.GetRequestID(ctx),
				"ip", c.RealIP(),
			)
			return err
		}
	}
}
