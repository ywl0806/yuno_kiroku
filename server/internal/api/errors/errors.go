package errors

import (
	"github.com/ywl0806/yuno_kiroku/internal/api/consts"
	"github.com/ywl0806/yuno_kiroku/internal/api/utils"
)

// Error types
type NotFoundError struct {
	Field string
}

func (e *NotFoundError) Error() string {
	return utils.GetMessage(consts.ErrNotFound, map[string]string{"field": e.Field})
}

func (e *NotFoundError) Is(target error) bool {
	_, ok := target.(*NotFoundError)
	return ok
}

type AlreadyExistsError struct {
	Field string
}

func (e *AlreadyExistsError) Error() string {
	return utils.GetMessage(consts.ErrAlreadyExists, map[string]string{"field": e.Field})
}

func (e *AlreadyExistsError) Is(target error) bool {
	_, ok := target.(*AlreadyExistsError)
	return ok
}

type InvalidError struct {
	Field string
}

func (e *InvalidError) Error() string {
	return utils.GetMessage(consts.ErrInvalid, map[string]string{"field": e.Field})
}

func (e *InvalidError) Is(target error) bool {
	_, ok := target.(*InvalidError)
	return ok
}

type RequiredError struct {
	Field string
}

func (e *RequiredError) Error() string {
	return utils.GetMessage(consts.ErrRequired, map[string]string{"field": e.Field})
}

func (e *RequiredError) Is(target error) bool {
	_, ok := target.(*RequiredError)
	return ok
}

type InternalError struct{}

func (e *InternalError) Error() string {
	return utils.GetMessage(consts.ErrInternal, nil)
}

func (e *InternalError) Is(target error) bool {
	_, ok := target.(*InternalError)
	return ok
}

type DuplicateError struct {
	Field string
}

func (e *DuplicateError) Error() string {
	return utils.GetMessage(consts.ErrDuplicate, map[string]string{"field": e.Field})
}

func (e *DuplicateError) Is(target error) bool {
	_, ok := target.(*DuplicateError)
	return ok
}

// Helper functions
func NewNotFoundError(field string) error {
	return &NotFoundError{Field: field}
}

func NewAlreadyExistsError(field string) error {
	return &AlreadyExistsError{Field: field}
}

func NewInvalidError(field string) error {
	return &InvalidError{Field: field}
}

func NewRequiredError(field string) error {
	return &RequiredError{Field: field}
}

func NewInternalError() error {
	return &InternalError{}
}

func NewDuplicateError(field string) error {
	return &DuplicateError{Field: field}
}
