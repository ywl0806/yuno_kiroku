package db

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"
)

// DB 커넥션 초기화
func Init(logQueries bool) (DBTX, error) {
	conn, err := sql.Open("pgx", viper.GetString("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	var dbtx DBTX = conn
	if logQueries {
		dbtx = NewQueryLogger(conn, true)
	}
	return dbtx, nil
}
