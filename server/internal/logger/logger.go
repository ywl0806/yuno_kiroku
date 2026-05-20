package logger

import (
	"log/slog"
	"os"
)

// Init은 앱 시작 시 한 번 호출한다.
// env, service 필드가 이후 모든 로그에 자동 포함된다.
func Init(env, service string) {
	var handler slog.Handler
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	if env == "local" || env == "dev" {
		opts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}
	slog.SetDefault(slog.New(handler).With(
		"env", env,
		"service", service,
	))
}
