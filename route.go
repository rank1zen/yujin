package main

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rank1zen/yujin/postgresql"
)

func RegisterRoutes(e *echo.Echo, pool *pgxpool.Pool) {
	v1 := e.Group("/v1", CheckHealth(pool))
	q := postgresql.NewQueries(pool)

	summonerv4 := v1.Group("/summoner/v4")
	{
		summonerv4.GET("/summoner/:uuid", GetSummoner(q))
		summonerv4.POST("/summoner", PostSummoner(q))
		summonerv4.DELETE("/summoner/:uuid", DeleteSummoner(q))

		summonerv4.GET("/by/puuid/:puuid", GetSummonerByPuuid(q))
		summonerv4.GET("/by/puuid/:puuid/recent", GetSummonerByPuuidRecent(q))
		summonerv4.GET("/by/puuid/:puuid/count", GetSummonerByPuuidCount(q))

		summonerv4.GET("/by/name/:name", GetSummonerByName(q))
		summonerv4.GET("/by/name/:name/recent", GetSummonerByNameRecent(q))
		summonerv4.GET("/by/name/:name/count", GetSummonerByNameCount(q))
	}
}
