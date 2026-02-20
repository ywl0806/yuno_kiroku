package db

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/utils"
)

// queryLogger는 DBTX를 래핑하고 쿼리를 로깅
type queryLogger struct {
	inner   DBTX
	enabled bool
}

// NewQueryLogger는 쿼리 로깅을 활성화할 때 DBTX를 반환
func NewQueryLogger(inner DBTX, enabled bool) DBTX {
	if !enabled {
		return inner
	}
	return &queryLogger{inner: inner, enabled: true}
}

// 쿼리와 리퀘스트ID를 로깅
func (q *queryLogger) logQuery(ctx context.Context, query string, args ...interface{}) {
	query = strings.TrimSpace(query)

	requestId := "unknown"
	if ec, ok := ctx.(echo.Context); ok {
		requestId = utils.GetRequestID(ec)
	}
	log.Printf("[DB] request_id=%s | query: %s | args_count=%d", requestId, query, len(args))
}

func (q *queryLogger) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if q.enabled {
		q.logQuery(ctx, query, args...)
	}
	return q.inner.ExecContext(ctx, query, args...)
}

func (q *queryLogger) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	if q.enabled {
		q.logQuery(ctx, query)
	}
	return q.inner.PrepareContext(ctx, query)
}

func (q *queryLogger) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if q.enabled {
		q.logQuery(ctx, query, args...)
	}
	return q.inner.QueryContext(ctx, query, args...)
}

func (q *queryLogger) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if q.enabled {
		q.logQuery(ctx, query, args...)
	}
	return q.inner.QueryRowContext(ctx, query, args...)
}
