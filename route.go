package main

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, p *pgxpool.Pool) {
	v1 := e.Group("/v1")

	m := []echo.MiddlewareFunc{
		MiddleDbHealth(p),
	}

	summoner := v1.Group("/summoner")
	summoner.GET("/puuid/by/:puuid", HandleGetSummonerRecordsByPuuid(), m...)
	summoner.GET("/puuid/count/:puuid", HandleGetSummonerRecordCountByPuuid(), m...)
	summoner.GET("/name/:name", HandleGetSummonerRecordsByName(), m...)
	summoner.POST("/", HandlePostSummonerRecord(), m...)

	soloq := v1.Group("/soloq")
	soloq.GET("/id/by/:id", HandleGetSoloqRecordById(), m...)
	soloq.GET("/name/by/:name", HandleGetSoloqRecordByName(), m...)
	soloq.POST("/", HandlePostSoloqRecord(), m...)

	match := v1.Group("/match")
	match.GET("/id/:id", HandleGetMatch(), m...)
	match.POST("/", HandlePostMatch(), m...)
}
