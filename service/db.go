package service

import (
	"fmt"
	"log"

	"github.com/crowmw/risiti/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	DB_NAME = "./data/risiti.db"
)

func DefaultDB() *gorm.DB {
	db, err := getConnection(DB_NAME)
	if err != nil {
		log.Fatal("Cannot get Sqlite DB Connection", err)
	}

	if err := migrateModels(db); err != nil {
		log.Fatal(err)
	}

	return db
}

func getConnection(dbName string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("🔥 failed to connect to the database: %s", err)
	}

	return db, nil
}

func migrateModels(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.User{}, &model.Receipt{}); err != nil {
		log.Fatal("Cannot create migratons", err)
	}

	return nil
}
