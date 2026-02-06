package errs

import (
	"runtime"
)

// Error types (재사용 가능, API/메시지 포맷 무관)

// 권한 오류 공통 타입
type AuthorizationError struct {
	Message string
}

func (e *AuthorizationError) Error() string {
	return e.Message
}

func (e *AuthorizationError) Is(target error) bool {
	_, ok := target.(*AuthorizationError)
	return ok
}

// 인증 오류 공통 타입
type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

func (e *UnauthorizedError) Is(target error) bool {
	_, ok := target.(*UnauthorizedError)
	return ok
}

// 내부 오류 공통 타입
type InternalError struct {
	Message string
	File    string
	Line    int
}

func (e *InternalError) Error() string {
	return e.Message
}

func (e *InternalError) Is(target error) bool {
	_, ok := target.(*InternalError)
	return ok
}

// Helper functions

// 권한 오류 생성
func NewAuthorizationError(message string) error {
	return &AuthorizationError{Message: message}
}

// 인증 오류 생성
func NewUnauthorizedError(message string) error {
	return &UnauthorizedError{Message: message}
}

// 내부 오류 생성
func NewInternalError(err error) error {
	_, file, line, _ := runtime.Caller(1)
	return &InternalError{
		Message: err.Error(),
		File:    file,
		Line:    line,
	}
}
