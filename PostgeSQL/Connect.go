package Postgesql

import (
	"database/sql"
	"fmt"
	"telegram_test_bot/Config"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	config := Config.GetConfig()
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", config.PostgreSQL.User, config.PostgreSQL.Password, config.PostgreSQL.DbName, config.PostgreSQL.SslMode)
	db, err := sql.Open("postgres", connStr)
	CheckError(err)

	return db, err
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
