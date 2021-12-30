package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var Conn *sql.DB = nil

func Create() error {
	if Conn != nil {
		return nil
	}

	db, err := sql.Open("sqlite3", "./schism.sqlite")
	if err != nil {
		Conn = db
	}
	return err
}
