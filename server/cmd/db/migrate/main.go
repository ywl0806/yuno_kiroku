package main

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/pkg/setting"
)

func main() {
	setting.SettingEnv()

	migrationsPath, err := filepath.Abs("db/migrations")
	if err != nil {
		slog.Error("마이그레이션 경로 오류", "error", err)
		os.Exit(1)
	}

	slog.Info("마이그레이션 시작", "path", migrationsPath, "database_url", viper.GetString("DATABASE_URL"))

	m, err := migrate.New(
		"file://"+migrationsPath,
		viper.GetString("DATABASE_URL"),
	)
	if err != nil {
		slog.Error("migrate 초기화 실패", "error", err)
		os.Exit(1)
	}

	if err = m.Up(); err != nil {
		slog.Error("마이그레이션 실패", "error", err)
		os.Exit(1)
	}

	slog.Info("마이그레이션 성공")

	if _, err = m.Close(); err != nil {
		slog.Error("migrate 종료 실패", "error", err)
		os.Exit(1)
	}

	slog.Info("마이그레이션 종료")
}
