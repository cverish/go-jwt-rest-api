package database

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB holds an instance of a gorm.DB object.
var DB *gorm.DB

// Database is a struct which points to a specific instance
// of a gorm.DB object.
type Database struct {
	DB *gorm.DB
}

// NewDatabase takes in a DSN string and returns an instance of Database.
// It opens the database via gorm with the postgres dialect, configures
// the connection pool, and pings the database.
func NewDatabase(dsn string) (*Database, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}
