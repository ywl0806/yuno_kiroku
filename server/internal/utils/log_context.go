package utils

import (
	"context"
	"log/slog"

	"github.com/ywl0806/yuno_kiroku/internal/consts"
)

type contextKey string

const (
	logUserIDKey   contextKey = "log_user_id"
	logFamilyIDKey contextKey = "log_family_id"
)

func WithLogUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, logUserIDKey, id)
}

func WithLogFamilyID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, logFamilyIDKey, id)
}

// LogAttrs는 context에서 공통 로그 필드(request_id, user_id, family_id)를 꺼내 반환한다.
// handler에서 service로 context를 전달할 때 자동으로 포함된다.
func LogAttrs(ctx context.Context) []any {
	var attrs []any
	if v := GetRequestID(ctx); v != "" {
		attrs = append(attrs, slog.String("request_id", v))
	}
	if v, ok := ctx.Value(logUserIDKey).(string); ok && v != "" {
		attrs = append(attrs, slog.String("user_id", v))
	}
	if v, ok := ctx.Value(logFamilyIDKey).(string); ok && v != "" {
		attrs = append(attrs, slog.String("family_id", v))
	}
	return attrs
}

// InjectAuthToContext는 인증된 사용자 정보를 로깅용으로 context에 심는다.
// Guard 미들웨어 이후 핸들러에서 호출한다.
func InjectAuthToContext(ctx context.Context, userID, familyID string) context.Context {
	ctx = context.WithValue(ctx, consts.RequestIDKey, GetRequestID(ctx))
	ctx = WithLogUserID(ctx, userID)
	ctx = WithLogFamilyID(ctx, familyID)
	return ctx
}
