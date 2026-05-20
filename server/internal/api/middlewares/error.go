package middlewares

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/apperr"
	"github.com/ywl0806/yuno_kiroku/internal/i18n"
	"github.com/ywl0806/yuno_kiroku/internal/utils"
	"github.com/ywl0806/yuno_kiroku/pkg/errs"
)

type ErrorHandler struct {
	log *slog.Logger
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{
		log: slog.Default().With("layer", "middleware", "component", "error"),
	}
}

func (e *ErrorHandler) Handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		return e.HandleServiceError(c, err)
	}
}

// 에러를 처리하고 HTTP 상태 코드를 반환합니다.
func (e *ErrorHandler) HandleServiceError(c echo.Context, err error) error {
	if err == nil {
		return nil
	}

	locale := GetLocale(c)

	// Echo의 HTTPError는 이미 처리된 에러이므로 그대로 반환
	if httpErr, ok := err.(*echo.HTTPError); ok {
		return httpErr
	}

	// 어플리케이션 에러 처리 (i18n 번역 적용)
	var appErrors *apperr.AppErrors
	if errors.As(err, &appErrors) {
		msg := i18n.T(locale, appErrors.Message, appErrors.TemplateData)

		switch appErrors.Code {
		case apperr.Unauthorized:
			return echo.NewHTTPError(http.StatusUnauthorized, msg)
		case apperr.Forbidden:
			return echo.NewHTTPError(http.StatusForbidden, msg)
		case apperr.Conflict:
			return echo.NewHTTPError(http.StatusConflict, msg)
		case apperr.Validation:
			return echo.NewHTTPError(http.StatusBadRequest, msg)
		case apperr.NotFound:
			return echo.NewHTTPError(http.StatusNotFound, msg)
		case apperr.BadRequest:
			return echo.NewHTTPError(http.StatusBadRequest, msg)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, msg)
		}
	}

	// 내부 에러 처리
	var internalError *errs.InternalError
	if errors.As(err, &internalError) {
		requestId := utils.GetRequestID(c.Request().Context())
		if requestId == "" {
			requestId = "unknown"
		}
		e.log.ErrorContext(c.Request().Context(), "internal error",
			"request_id", requestId,
			"file", internalError.File,
			"line", internalError.Line,
			"message", internalError.Message,
		)
		return echo.NewHTTPError(http.StatusInternalServerError, internalError.Message)
	}

	// 인증 에러 처리
	var unauthorizedError *errs.UnauthorizedError
	if errors.As(err, &unauthorizedError) {
		return echo.NewHTTPError(http.StatusUnauthorized, unauthorizedError.Message)
	}

	// 권한 에러 처리
	var authorizationError *errs.AuthorizationError
	if errors.As(err, &authorizationError) {
		return echo.NewHTTPError(http.StatusForbidden, authorizationError.Message)
	}

	// 미분류 에러
	return echo.NewHTTPError(http.StatusInternalServerError, i18n.T(locale, "message.internal", nil))
}
