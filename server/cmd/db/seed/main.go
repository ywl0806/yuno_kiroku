package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/pkg/setting"
	"github.com/ywl0806/yuno_kiroku/pkg/utils"
)

func main() {
	setting.SettingEnv()

	databaseURL := viper.GetString("DATABASE_URL")
	env := viper.GetString("APP_ENV")
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		slog.Error("DB 연결 실패", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	seedFilePath := fmt.Sprintf("db/seed.%s.sql", env)
	seedSQL, err := os.ReadFile(seedFilePath)
	if err != nil {
		slog.Error("시드 파일 읽기 실패", "path", seedFilePath, "error", err)
		os.Exit(1)
	}

	slog.Info("데이터 시드 중", "env", env)
	if _, err = db.Exec(string(seedSQL)); err != nil {
		slog.Error("시드 SQL 실행 실패", "error", err)
		os.Exit(1)
	}

	users, err := db.Query("SELECT id, password FROM users")
	if err != nil {
		slog.Error("users 조회 실패", "error", err)
		os.Exit(1)
	}
	defer users.Close()

	for users.Next() {
		var id string
		var password string
		if err = users.Scan(&id, &password); err != nil {
			slog.Error("users scan 실패", "error", err)
			os.Exit(1)
		}
		hashPassword, err := utils.HashPassword(password)
		if err != nil {
			slog.Error("비밀번호 해시 실패", "error", err)
			os.Exit(1)
		}
		if _, err = db.Exec("UPDATE users SET password = $1 WHERE id = $2", hashPassword, id); err != nil {
			slog.Error("비밀번호 업데이트 실패", "error", err)
			os.Exit(1)
		}
	}

	slog.Info("데이터 시드 완료")
}
