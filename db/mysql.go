package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbPool *sql.DB
	dbUrl  = "root:root@/echo_demo?charset=utf8&parseTime=True&loc=Local"
)

func ConnInit() error {
	db, err := sql.Open("mysql", dbUrl)
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
