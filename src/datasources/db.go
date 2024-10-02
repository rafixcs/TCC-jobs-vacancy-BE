package datasources

import (
	"database/sql"
	"fmt"
)

func OpenDb() (*sql.DB, error) {
	dsn := fmt.Sprintf(`host=%s user=%s password=%s port=%s dbname=%s sslmode=disable`, "localhost", "root", "root", "1234", "root")
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
