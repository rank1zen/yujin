package main

import (
	"log"
	"os"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/gin-gonic/gin"
	"github.com/rank1zen/yujin/internal/postgresql"
	"github.com/rank1zen/yujin/internal/rest"
	"github.com/rank1zen/yujin/internal/riot"
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

	client := golio.NewClient(os.Getenv("RIOT_API_KEY"), golio.WithRegion(api.RegionNorthAmerica))
	riot := riot.NewSummonerQ(client.Riot.LoL.Summoner)

	da := postgresql.NewSummonerDA(db)
	svc := service.NewSummonerService(da, riot)

	hdl := rest.NewSummonerHandler(svc)
	hdl.Register(r)

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}
