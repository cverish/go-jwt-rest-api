package testutils

import (
	"errors"
	"os"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cverish/go-jwt-rest-api/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetMockGormDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, errors.New("error initializing sqlmock db")
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, errors.New("error initializing gorm db")
	}

	return gormDB, mock, nil
}

func GetMockConfig() (*config.Config, error) {
	f, err := os.CreateTemp(".", ".env.foo")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())

	cfg, err := config.Load(f.Name())
	return cfg, err
}
