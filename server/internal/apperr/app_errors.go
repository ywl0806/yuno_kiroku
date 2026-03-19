package apperr

type ErrorCode string

// 에러 코드
const (
	Unauthorized     ErrorCode = "UNAUTHORIZED"
	Forbidden        ErrorCode = "FORBIDDEN"
	Conflict         ErrorCode = "CONFLICT"
	Validation       ErrorCode = "VALIDATION"
	Internal         ErrorCode = "INTERNAL"
	NotFound         ErrorCode = "NOT_FOUND"
	Duplicate        ErrorCode = "DUPLICATE"
	InvalidReference ErrorCode = "INVALID_REFERENCE"
	BadRequest       ErrorCode = "BAD_REQUEST"
)

// 어플리케이션 에러 타입
// Code: 에러 코드
// Message: i18n 메시지 키 또는 메시지 문자열
// TemplateData: Message가 키일 때 템플릿 치환 데이터 (예: field.group)
// Err: 에러 객체
type AppErrors struct {
	Code         ErrorCode
	Message      string
	TemplateData map[string]string
	Err          error
}

func (e *AppErrors) Error() string {
	return e.Message
}

func (e *AppErrors) Is(target error) bool {
	_, ok := target.(*AppErrors)
	return ok
}

func IsAppError(err error, code ErrorCode) bool {
	appErr := &AppErrors{Code: code}
	if appErr.Is(err) {
		return appErr.Code == err.(*AppErrors).Code
	}
	return false
}

// 인증 에러 생성
func NewUnauthorizedError(message string, templateData map[string]string) error {
	return &AppErrors{Code: Unauthorized, Message: message}
}

// 권한 에러 생성
func NewForbiddenError(message string, templateData map[string]string) error {
	return &AppErrors{Code: Forbidden, Message: message, TemplateData: templateData}
}

// 충돌 에러 생성
func NewConflictError(message string, templateData map[string]string) error {
	return &AppErrors{Code: Conflict, Message: message, TemplateData: templateData}
}

// 유효성 에러 생성
func NewValidationError(message string, templateData map[string]string) error {
	return &AppErrors{Code: Validation, Message: message, TemplateData: templateData}
}

// 존재하지 않는 에러 생성
func NewNotFoundError(message string, templateData map[string]string) error {
	return &AppErrors{Code: NotFound, Message: message, TemplateData: templateData}
}

// 중복 에러 생성
func NewDuplicateError(message string, templateData map[string]string) error {
	return &AppErrors{Code: Duplicate, Message: message, TemplateData: templateData}
}

// 참조 오류 생성
func NewInvalidReferenceError(message string, templateData map[string]string) error {
	return &AppErrors{Code: InvalidReference, Message: message, TemplateData: templateData}
}

// 잘못된 요청 에러 생성
func NewBadRequestError(message string, templateData map[string]string) error {
	return &AppErrors{Code: BadRequest, Message: message, TemplateData: templateData}
}

// i18n 메시지 키 + 템플릿 데이터로 에러 생성
func NewAppErrorWithData(code ErrorCode, messageKey string, templateData map[string]string) error {
	return &AppErrors{Code: code, Message: messageKey, TemplateData: templateData}
}
