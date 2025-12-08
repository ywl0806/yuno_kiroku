package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	apiErrors "github.com/ywl0806/yuno_kiroku/internal/api/errors"
)

// HandleServiceError converts service layer errors to appropriate HTTP errors
func HandleServiceError(err error) error {
	if err == nil {
		return nil
	}

	// Check for specific error types
	if err == apiErrors.ErrNotFound ||
		err == apiErrors.ErrPhotoNotFound ||
		err == apiErrors.ErrUserNotFound ||
		err == apiErrors.ErrIdentityNotFound ||
		err == apiErrors.ErrGroupNotFound ||
		err == apiErrors.ErrClanGroupNotFound {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	if err == apiErrors.ErrUsernameAlreadyExists {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}

	// For other errors, return as-is (will be handled as 500 by echo)
	return err
}

