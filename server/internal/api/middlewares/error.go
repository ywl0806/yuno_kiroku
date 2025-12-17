package middlewares

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	customErrors "github.com/ywl0806/yuno_kiroku/internal/api/errors"
)

type ErrorHandler struct {
	notFoundErr      *customErrors.NotFoundError
	alreadyExistsErr *customErrors.AlreadyExistsError
	invalidErr       *customErrors.InvalidError
	requiredErr      *customErrors.RequiredError
	internalErr      *customErrors.InternalError
	duplicateErr     *customErrors.DuplicateError
}

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

	// 커스텀 에러 타입 체크
	switch {
	case errors.As(err, &e.notFoundErr):
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	case errors.As(err, &e.alreadyExistsErr):
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	case errors.As(err, &e.invalidErr), errors.As(err, &e.requiredErr):
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	case errors.As(err, &e.internalErr):
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	case errors.As(err, &e.duplicateErr):
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	default:
		// 알 수 없는 에러는 500으로 처리
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
}
