package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/pkg/setting"
)

func main() {
	setting.SettingEnv()

	databaseURL := viper.GetString("DATABASE_URL")

	env := viper.GetString("APP_ENV")
	if !(env == "local" || env == "dev") {
		log.Fatal("This command is only available in development and local environment")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("스키마 삭제중...")
	_, err = db.Exec("DROP SCHEMA public CASCADE;")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("DROP TABLE IF EXISTS schema_migrations;")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("스키마 재생성중...")
	_, err = db.Exec("CREATE SCHEMA public;")

	if err != nil {
		log.Fatal(err)
	}
	log.Println("스키마 삭제 및 재생성 완료")

}
