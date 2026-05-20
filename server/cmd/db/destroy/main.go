package main

import (
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/pkg/setting"
)

func main() {
	setting.SettingEnv()

	databaseURL := viper.GetString("DATABASE_URL")

	env := viper.GetString("APP_ENV")
	if !(env == "local" || env == "dev") {
		slog.Error("This command is only available in development and local environment", "env", env)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		slog.Error("DB 연결 실패", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("스키마 삭제 중")
	if _, err = db.Exec("DROP SCHEMA IF EXISTS public CASCADE;"); err != nil {
		slog.Error("스키마 삭제 실패", "error", err)
		os.Exit(1)
	}
	if _, err = db.Exec("DROP TABLE IF EXISTS schema_migrations;"); err != nil {
		slog.Error("schema_migrations 삭제 실패", "error", err)
		os.Exit(1)
	}

	slog.Info("스키마 재생성 중")
	if _, err = db.Exec("CREATE SCHEMA public;"); err != nil {
		slog.Error("스키마 재생성 실패", "error", err)
		os.Exit(1)
	}

	slog.Info("스키마 삭제 및 재생성 완료")
}
