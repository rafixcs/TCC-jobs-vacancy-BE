package datasources

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	config "github.com/rafixcs/tcc-job-vacancy/src/configuration"
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
	dsn := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`, config.DB_HOST, config.DB_PORT, config.DB_USER, config.DB_PASSWORD, config.DB_NAME)
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
