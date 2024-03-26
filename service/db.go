package service

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DB_NAME = "./data/risiti.db"
)

func NewDB() *sql.DB {
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

	// Init SQLite3 database
	db, err = sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, fmt.Errorf("🔥 failed to connect to the database: %s", err)
	}

	return db, nil
}

func createMigrations(db *sql.DB) error {
	stmt := `CREATE TABLE IF NOT EXISTS user (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL
	);`

	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}

	stmt = `CREATE TABLE IF NOT EXISTS receipt (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL UNIQUE,
		filename VARCHAR(255) NOT NULL UNIQUE,
		description VARCHAR(255) NOT NULL,
		date DATETIME NOT NULL,
		created_by INTEGER,
		FOREIGN KEY(created_by) REFERENCES user(id)
	);`

	_, err = db.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}
