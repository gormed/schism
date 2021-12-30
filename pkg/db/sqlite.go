package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type CRUD interface {
	Create()
	Read()
	Update()
	Delete()
}

type Sqlite struct {
	Conn *sql.DB
}

func NewSqlite() *Sqlite {
	return &Sqlite{}
}

func (s *Sqlite) Create() error {
	if s.Conn != nil {
		return nil
	}

	db, err := sql.Open("sqlite3", "./schism.sqlite")
	if err != nil {
		s.Conn = db
	}
	return err
}

func (s *Sqlite) Close() {
	defer s.Conn.Close()
}

func (s *Sqlite) Insert(table string, fields map[string]interface{}) (*sql.Result, error) {
	if s.Conn == nil {
		return nil, fmt.Errorf("database connection is missing, use s.Create() before any other function")
	}
	if len(table) == 0 {
		return nil, fmt.Errorf("no valid table name given '%s'", table)
	}
	var keys = ""
	var values = ""
	var valuesArray = []interface{}{}
	for key, value := range fields {
		keys += key + ", "
		values += "?, "
		valuesArray = append(valuesArray, value)
	}
	sqlStmt := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, table, keys, values)
	stmt, err := s.Conn.Prepare(sqlStmt)
	if err != nil {
		return nil, err
	}
	result, err := stmt.Exec(valuesArray...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *Sqlite) Prepare(queryStmt string) (*sql.Stmt, error) {
	if s.Conn == nil {
		return nil, fmt.Errorf("database connection is missing, use s.Create() before any other function")
	}
	stmt, err := s.Conn.Prepare(queryStmt)
	if err != nil {
		return nil, err
	}
	return stmt, nil
}
