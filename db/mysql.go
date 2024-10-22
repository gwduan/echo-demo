package db

import (
	"database/sql"
	"echo-demo/config"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

var (
	ErrDupRows  = errors.New("DB: Duplicate")
	ErrNotFound = errors.New("DB: Not Found")
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
