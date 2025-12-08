package errors

import "errors"

// Common application errors
var (
	ErrNotFound              = errors.New("resource not found")
	ErrPhotoNotFound         = errors.New("photo not found")
	ErrUserNotFound          = errors.New("user not found")
	ErrIdentityNotFound      = errors.New("identity not found")
	ErrGroupNotFound         = errors.New("group not found")
	ErrClanGroupNotFound     = errors.New("clan group not found")
	ErrUsernameAlreadyExists = errors.New("username already exists")
)
