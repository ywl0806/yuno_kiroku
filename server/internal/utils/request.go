package utils

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/consts"
)

// echo context에서 request_id를 가져온다.
func GetRequestID(c context.Context) string {
	v, ok := c.Value(consts.RequestIDKey).(string)
	if !ok {
		return ""
	}
	return v
}
