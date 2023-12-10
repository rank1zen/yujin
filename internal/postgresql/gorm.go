package postgresql

import (
	"gorm.io/gorm"
)

func Setup(q *gorm.DB) {
	q.AutoMigrate(&Summoner{})
}
