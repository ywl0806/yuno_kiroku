package main

import (
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/pkg/setting"
)

func main() {
	setting.SettingEnv()

	// 마이그레이션 경로 설정
	migrationsPath, err := filepath.Abs("db/migrations")
	if err != nil {
		log.Fatal("마이그레이션 경로 오류: ", err)
	}

	log.Printf("마이그레이션 경로: %s", migrationsPath)
	log.Printf("DATABASE_URL: %s", viper.GetString("DATABASE_URL"))

	m, err := migrate.New(
		"file://"+migrationsPath,
		viper.GetString("DATABASE_URL"),
	)

	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("마이그레이션 성공")

	_, err = m.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("마이그레이션 종료")
}
