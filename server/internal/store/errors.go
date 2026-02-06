package store

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/ywl0806/yuno_kiroku/internal/apperr"
	"github.com/ywl0806/yuno_kiroku/pkg/errs"
)

// DB 에러를 어플리케이션 에러로 변환
func mapDBError(err error, notFoundResource ...string) error {
	if err == nil {
		return nil
	}

	resource := "resource"
	if len(notFoundResource) > 0 && notFoundResource[0] != "" {
		resource = notFoundResource[0]
	}

	if errors.Is(err, sql.ErrNoRows) {
		return apperr.NewNotFoundError(resource)
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505": // unique_violation
			return apperr.NewConflictError("duplicate")
		case "23503": // foreign_key_violation
			return apperr.NewInvalidReferenceError(resource)
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
