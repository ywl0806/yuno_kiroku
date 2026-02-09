package store

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/ywl0806/yuno_kiroku/internal/apperr"
	"github.com/ywl0806/yuno_kiroku/pkg/errs"
)

// DB 에러를 어플리케이션 에러로 변환
// notFoundResource: i18n 필드 키 (예: "field.group", "field.user"). 비어 있으면 "field.resource" 사용.
func mapDBError(err error, notFoundResource ...string) error {
	if err == nil {
		return nil
	}

	fieldKey := "field.resource"
	if len(notFoundResource) > 0 && notFoundResource[0] != "" {
		fieldKey = notFoundResource[0]
	}

	if errors.Is(err, sql.ErrNoRows) {
		return apperr.NewNotFoundError("message.not_found", map[string]string{"field": fieldKey})
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505": // unique_violation
			return apperr.NewConflictError("message.conflict.duplicate", map[string]string{"field": fieldKey})
		case "23503": // foreign_key_violation
			return apperr.NewInvalidReferenceError("message.invalid_reference", map[string]string{"field": fieldKey})
		}
	}

	return errs.NewInternalError(err)
}

// wrapErr DB 호출 결과 (T, error)에 대해 err가 있으면 mapDBError로 감싸서 반환한다.
// notFoundResource는 mapDBError에 전달할 리소스 이름(선택).
func wrapErr[T any](v T, err error, notFoundResource ...string) (T, error) {
	if err != nil {
		var zero T
		return zero, mapDBError(err, notFoundResource...)
	}
	return v, nil
}
