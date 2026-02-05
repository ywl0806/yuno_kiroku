package appErrors

type ErrorCode string

// 에러 코드
const (
	Unauthorized ErrorCode = "UNAUTHORIZED"
	Forbidden    ErrorCode = "FORBIDDEN"
	Conflict     ErrorCode = "CONFLICT"
	Validation   ErrorCode = "VALIDATION"
	Internal     ErrorCode = "INTERNAL"
	NotFound     ErrorCode = "NOT_FOUND"
	Duplicate    ErrorCode = "DUPLICATE"
)

// 어플리케이션 에러 타입
// Code: 에러 코드
// Message: 에러 메시지
// Err: 에러 객체
type AppErrors struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *AppErrors) Error() string {
	return e.Message
}

func (e *AppErrors) Is(target error) bool {
	_, ok := target.(*AppErrors)
	return ok
}

// 인증 에러 생성
func NewUnauthorizedError(message string) error {
	return &AppErrors{Code: Unauthorized, Message: message}
}

// 권한 에러 생성
func NewForbiddenError(message string) error {
	return &AppErrors{Code: Forbidden, Message: message}
}

// 충돌 에러 생성
func NewConflictError(message string) error {
	return &AppErrors{Code: Conflict, Message: message}
}

// 유효성 에러 생성
func NewValidationError(message string) error {
	return &AppErrors{Code: Validation, Message: message}
}

// 존재하지 않는 에러 생성
func NewNotFoundError(message string) error {
	return &AppErrors{Code: NotFound, Message: message}
}

// 중복 에러 생성
func NewDuplicateError(message string) error {
	return &AppErrors{Code: Duplicate, Message: message}
}
