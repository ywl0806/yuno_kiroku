package appErrors

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/ywl0806/yuno_kiroku/pkg/commonErrors"
)

// ClassifyDBError maps DB/driver errors to app or common errors.
// - sql.ErrNoRows → NotFound (resource name from notFoundResource, or "resource")
// - PostgreSQL 23505 (unique_violation) → Conflict
// - PostgreSQL 23503 (foreign_key_violation) → NotFound("referenced resource")
// - other → InternalError (logged in middleware, 500)
func ClassifyDBError(err error, notFoundResource ...string) error {
	if err == nil {
		return nil
	}

	resource := "resource"
	if len(notFoundResource) > 0 && notFoundResource[0] != "" {
		resource = notFoundResource[0]
	}

	if errors.Is(err, sql.ErrNoRows) {
		return NewNotFoundError(resource)
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505": // unique_violation
			return NewConflictError("duplicate")
		case "23503": // foreign_key_violation
			return NewNotFoundError("referenced " + resource)
		}
	}

	return commonErrors.NewInternalError(err)
}
