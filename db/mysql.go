package db

import (
	"database/sql"
	"echo-demo/config"

	_ "github.com/go-sql-driver/mysql"
)

var dbPool *sql.DB

func ConnInit() error {
	db, err := sql.Open(config.DbName(), config.DbUrl())
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	dbPool = db

	return nil
}

func Conn() *sql.DB {
	return dbPool
}
