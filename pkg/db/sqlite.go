package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var SqliteDateLayout = "2006-01-02T15:04:05-0700"

type Sqlite struct {
	conn *sql.DB
}

func NewSqlite() *Sqlite {
	return &Sqlite{conn: nil}
}

func (s *Sqlite) Create() error {
	if s.conn != nil {
		return nil
	}

	db, err := sql.Open("sqlite3", "/db/schism.sqlite")
	if err != nil {
		return err
	}

	s.conn = db

	err = s.setupDatabase()
	if err != nil {
		return err
	}
	return nil
}

func (s *Sqlite) Close() {
	defer s.conn.Close()
}

func (s *Sqlite) Prepare(queryStmt string) (*sql.Stmt, error) {
	if s.conn == nil {
		return nil, fmt.Errorf("database connection is missing, use s.Create() before any other function")
	}
	stmt, err := s.conn.Prepare(queryStmt)
	if err != nil {
		return nil, err
	}
	return stmt, nil
}
