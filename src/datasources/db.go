package datasources

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type IDatabasePsql interface {
	Open()
	Close()
	GetDB() *sql.DB
	GetError() error
}

type DatabasePsql struct {
	DB    *sql.DB
	Error error
}

func (dbpsql *DatabasePsql) Open() {
	dsn := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`, "localhost", "1234", "root", "root", "root")
	dbpsql.DB, dbpsql.Error = sql.Open("postgres", dsn)
	if dbpsql.Error != nil {
		return
	}

	dbpsql.Error = dbpsql.DB.Ping()
	if dbpsql.Error != nil {
		return
	}
}

func (dbpsql *DatabasePsql) Close() {
	dbpsql.DB.Close()
}

func (dbpsql *DatabasePsql) GetDB() *sql.DB {
	return dbpsql.DB
}

func (dbpsql *DatabasePsql) GetError() error {
	return dbpsql.Error
}
