package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rank1zen/yujin/internal/postgresql"
	"github.com/rank1zen/yujin/internal/rest"
	"github.com/rank1zen/yujin/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()


	dsn := "postgres://gordon:kop123456@localhost:5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	da := postgresql.NewSummonerDA(db)
	svc := service.NewSummonerService(da)
	hdl := rest.NewSummonerHandler(svc)
	hdl.Register(r)

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}
