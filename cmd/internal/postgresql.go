package internal

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgreSQL() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}
