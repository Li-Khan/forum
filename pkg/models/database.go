package models

import (
	"database/sql"
	"fmt"
	"io/ioutil"
)

// OpenDB ...
func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("./pkg/models/sqlite/%s", dsn))
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	err = setup(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setup(db *sql.DB) error {
	query, err := ioutil.ReadFile("./pkg/models/sqlite/setup.sql")
	if err != nil {
		return fmt.Errorf("setup: %s", err)
	}
	if _, err := db.Exec(string(query)); err != nil {
		return fmt.Errorf("setup: %s", err)
	}
	return nil
}
