package db

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"
)

// DB 커넥션 초기화
// *sql.DB는 BeginTx 전용, DBTX는 쿼리 실행(QueryLogger 래핑) 전용
func Init(logQueries bool) (*sql.DB, DBTX, error) {
	conn, err := sql.Open("pgx", viper.GetString("DATABASE_URL"))
	if err != nil {
		return nil, nil, err
	}
	var dbtx DBTX = conn
	if logQueries {
		dbtx = NewQueryLogger(conn, true)
	}
	return conn, dbtx, nil
}
