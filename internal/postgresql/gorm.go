package postgresql

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/rank1zen/yujin/internal/postgresql/db"
)

func InitDB() *gorm.DB {
	da, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	da.AutoMigrate(&db.MSummoner{})

	da.Create(&db.MSummoner{
		PuuId: "asd",
		AccountId: "adsa",
		SummonerId: "ADada",
		Level: 123,
		ProfileIconId: 369,
		Name: "asd",
	})

	return da
}
