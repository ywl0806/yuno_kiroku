package middlewares

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/api/appErrors"
	"github.com/ywl0806/yuno_kiroku/pkg/commonErrors"
)

type ErrorHandler struct{}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func (e *ErrorHandler) Handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)

		return e.HandleServiceError(err)
	}
}

// 에러를 처리하고 HTTP 상태 코드를 반환합니다.
func (e *ErrorHandler) HandleServiceError(err error) error {
	if err == nil {
		return nil
	}

	// Echo의 HTTPError는 이미 처리된 에러이므로 그대로 반환
	if httpErr, ok := err.(*echo.HTTPError); ok {
		return httpErr
	}

	// 어플리케이션 에러 처리
	var appErr *appErrors.AppErrors
	if errors.As(err, &appErr) {
		switch appErr.Code {
		case appErrors.Unauthorized:
			return echo.NewHTTPError(http.StatusUnauthorized, appErr.Message)
		case appErrors.Forbidden:
			return echo.NewHTTPError(http.StatusForbidden, appErr.Message)
		case appErrors.Conflict:
			return echo.NewHTTPError(http.StatusConflict, appErr.Message)
		case appErrors.Validation:
			return echo.NewHTTPError(http.StatusBadRequest, appErr.Message)
		case appErrors.Internal:
			return echo.NewHTTPError(http.StatusInternalServerError, appErr.Message)
		case appErrors.NotFound:
			return echo.NewHTTPError(http.StatusNotFound, appErr.Message)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, appErr.Message)
		}
	}

	// 내부 에러 처리
	var internalError *commonErrors.InternalError
	if errors.As(err, &internalError) {
		log.Printf("%s:%d: %s\n", internalError.File, internalError.Line, internalError.Message)
		return echo.NewHTTPError(http.StatusInternalServerError, internalError.Message)
	}

	// 인증 에러 처리
	var unauthorizedError *commonErrors.UnauthorizedError
	if errors.As(err, &unauthorizedError) {
		return echo.NewHTTPError(http.StatusUnauthorized, unauthorizedError.Message)
	}

	// 권한 에러 처리
	var authorizationError *commonErrors.AuthorizationError
	if errors.As(err, &authorizationError) {
		return echo.NewHTTPError(http.StatusForbidden, authorizationError.Message)
	}

	// 미분류 에러
	return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
}
