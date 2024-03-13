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

type Store struct {
	db *sql.DB
}

func NewStore() *Store {
	db, err := getConnection(DB_NAME)
	if err != nil {
		log.Fatal("Cannot get Sqlite DB Connection", err)
	}

	if err := createMigrations(db); err != nil {
		log.Fatal("Cannot create migratons", err)
	}

	return &Store{
		db,
	}
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
		// log.Fatalf("🔥 failed to connect to the database: %s", err.Error())
		return nil, fmt.Errorf("🔥 failed to connect to the database: %s", err)
	}

	slog.Info("🚀 Connected Successfully to the Database")

	return db, nil
}

func createMigrations(db *sql.DB) error {
	stmt := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		username VARCHAR(64) NOT NULL
	);`

	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}

	stmt = `CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_by INTEGER NOT NULL,
		title VARCHAR(64) NOT NULL,
		description VARCHAR(255) NULL,
		status BOOLEAN DEFAULT(FALSE),
		created_at DATETIME default CURRENT_TIMESTAMP,
		FOREIGN KEY(created_by) REFERENCES users(id)
	);`

	_, err = db.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}
