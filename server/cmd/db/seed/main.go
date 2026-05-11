package main

import (
	"database/sql"
	"fmt"
	"log"
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
		log.Fatal(err)
	}
	defer db.Close()

	seedFilePath := fmt.Sprintf("db/seed.%s.sql", env)
	seedSQL, err := os.ReadFile(seedFilePath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("데이터 시드중...")
	_, err = db.Exec(string(seedSQL))
	if err != nil {
		log.Fatal(err)
	}

	// 비밀번호 해시
	users, err := db.Query("SELECT id, password FROM users")

	if err != nil {
		log.Fatal(err)
	}
	defer users.Close()

	// 비밀번호 해시 업데이트
	for users.Next() {
		var id string
		var password string

		err = users.Scan(&id, &password)
		if err != nil {
			log.Fatal(err)
		}
		hashPassword, err := utils.HashPassword(password)
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Exec("UPDATE users SET password = $1 WHERE id = $2", hashPassword, id)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("데이터 시드 완료")
}
