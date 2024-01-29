package main

import (
	"github.com/labstack/echo/v4"
	"github.com/rank1zen/yujin/postgresql"
)

func RegisterRoutes(e *echo.Echo, q *postgresql.Queries) {
	v1 := e.Group("/v1")

	m := []echo.MiddlewareFunc{
		MiddleDbConn(p),
	}

	summonerv4 := v1.Group("/summonerv4")
	summonerv4.GET("/puuid/by/:puuid", HandleGetSummonerRecordsByPuuid(q), m...)
	summonerv4.GET("/puuid/count/:puuid", HandleGetSummonerRecordCountByPuuid(q), m...)
	summonerv4.GET("/name/by/:name", HandleGetSummonerRecordsByName(q), m...)
	summonerv4.POST("/", HandlePostSummonerRecord(q), m...)
	summonerv4.DELETE("/by", HandleDeleteSummonerRecord(q), m...)

	soloq := v1.Group("/soloq")
	soloq.GET("/id/by/:id", HandleGetSoloqRecordById(), m...)
	soloq.GET("/name/by/:name", HandleGetSoloqRecordByName(), m...)
	soloq.POST("/", HandlePostSoloqRecord(), m...)

	match := v1.Group("/match")
	match.GET("/id/:id", HandleGetMatch(), m...)
	match.POST("/", HandlePostMatch(), m...)
}
