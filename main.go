package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rank1zen/yujin/postgresql"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	pool, err := pgxpool.New(context.Background(), os.Getenv("YUJIN_PG_STRING"))
	if err != nil {
		log.Fatal("can't connect to database")
	}

	gc := golio.NewClient(os.Getenv("YUJIN_RIOT_API_KEY"), golio.WithRegion(api.RegionNorthAmerica))

	e.GET("/", HandleHome())
	v1 := e.Group("/v1", CheckHealth(pool))
	q := postgresql.NewQueries(pool)

	summonerv4 := v1.Group("/summoner/v4")
	{
		summonerv4.GET("/record/:uuid", GetSummoner(q))
		summonerv4.DELETE("/record/:uuid", DeleteSummoner(q))

		summonerv4.GET("/by-puuid/:puuid", GetSummonerByPuuid(q))
		summonerv4.GET("/by-puuid/:puuid/recent", GetSummonerByPuuidRecent(q))
		summonerv4.GET("/by-puuid/:puuid/count", GetSummonerByPuuidCount(q))

		summonerv4.GET("/by-name/:name", GetSummonerByName(q))
		summonerv4.GET("/by-name/:name/recent", GetSummonerByNameRecent(q))
		summonerv4.GET("/by-name/:name/count", GetSummonerByNameCount(q))
		summonerv4.POST("/by-name/:name/renew", PostSummonerByName(q, gc))
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		err := e.Start(fmt.Sprintf("0.0.0.0:%s", os.Getenv("YUJIN_PORT")))
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("shutting down server...")
		}
	}()

	<-ctx.Done()

	err = e.Shutdown(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
}
