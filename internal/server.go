package internal

import (
	"github.com/KnutZuidema/golio"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rank1zen/yujin/internal/postgresql"
)

func Routes(e *echo.Echo, pool *pgxpool.Pool, gc *golio.Client) {
	e.GET("/", HandleHome())
	v1 := e.Group("/v1")
	q := postgresql.NewQuery(pool)

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
}
