package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/pkg/utils"
	"github.com/ywl0806/yuno_kiroku/pkg/setting"
)

func main() {
	setting.SettingEnv()

	databaseURL := viper.GetString("DATABASE_URL")

	db, err := sql.Open("postgres", databaseURL)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	seedFilePath := "db/seed.sql"
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
	hashPassword, _ := utils.HashPassword("password")
	users, err := db.Query("SELECT id FROM users")

	if err != nil {
		log.Fatal(err)
	}
	defer users.Close()

	// 비밀번호 해시 업데이트
	for users.Next() {
		var id int

		err = users.Scan(&id)
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
