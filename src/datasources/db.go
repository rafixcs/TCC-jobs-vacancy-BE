package datasources

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func OpenDb() (*sql.DB, error) {
	dsn := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`, "localhost", "1234", "root", "root", "root")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
