package inter

import (
	"github.com/rank1zen/yujin/internal/postgresql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TODO: Configurations
func NewPostgreSQL() (*gorm.DB, error) {
	dsn := "postgres://gordon:kop123456@localhost:5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	postgresql.Setup(db)
	return db, err
}
