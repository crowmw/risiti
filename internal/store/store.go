package store

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DB_NAME = "risiti.db"
)

func NewStore() *sql.DB {
	db, err := getConnection(DB_NAME)
	if err != nil {
		log.Fatal("Cannot get Sqlite DB Connection", err)
	}

	if err := createMigrations(db); err != nil {
		log.Fatal("Cannot create migratons", err)
	}

	return db
}

func getConnection(dbName string) (*sql.DB, error) {
	var (
		err error
		db  *sql.DB
	)

	if db != nil {
		return db, nil
	}

	// Init SQLite3 database
	db, err = sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, fmt.Errorf("🔥 failed to connect to the database: %s", err)
	}

	slog.Info("🚀 Connected Successfully to the Database")

	return db, nil
}

func createMigrations(db *sql.DB) error {
	stmt := `CREATE TABLE IF NOT EXISTS receipt (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL UNIQUE,
		filename VARCHAR(255) NOT NULL UNIQUE,
		description VARCHAR(255) NOT NULL,
		date VARCHAR(255) NOT NULL
	);`

	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}
