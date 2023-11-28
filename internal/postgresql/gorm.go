package postgresql

import (
	"github.com/rank1zen/yujin/internal/postgresql/db"
	"gorm.io/gorm"
)

func Setup(q *gorm.DB) {
	q.AutoMigrate(&db.Summoner{})
}
