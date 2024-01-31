package main

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rank1zen/yujin/postgresql"
)

func HandleGetSummonerRecordsByPuuid(q *postgresql.Queries) echo.HandlerFunc {
	type PathParam struct {
		Puuid string `param:"puuid"`
	}

	return func(c echo.Context) error {
		var path PathParam
		if err := c.Bind(&path); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.String(http.StatusOK, "Not Implemented")
	}
}

func HandleGetSummonerRecordCountByPuuid(q *postgresql.Queries) echo.HandlerFunc {
	type PathParam struct {
		Puuid string `param:"puuid"`
	}

	return func(c echo.Context) error {
		var path PathParam
		if err := c.Bind(&path); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		count, err := q.CountSummonerRecordsByPuuid(ctx, path.Puuid)
		if err != nil {
			return echo.NewHTTPError(http.StatusServiceUnavailable, err.Error())
		}

		return c.JSON(http.StatusOK, count)
	}
}

func HandleGetSummonerRecordsByName(q *postgresql.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Not Implemented")
	}
}

func HandlePostSummonerRecord(q *postgresql.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		var body SummonerRecordBody
		if err := c.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		id, err := q.InsertSummonerRecord(ctx, &postgresql.SummonerRecord{})
		if err != nil {
			return echo.NewHTTPError(http.StatusServiceUnavailable, err.Error())
		}

		return c.JSON(http.StatusOK, id)
	}
}

func HandleDeleteSummonerRecord(q *postgresql.Queries) echo.HandlerFunc {
	type QueryParam struct {
		Id string `query:"id"`
	}

	return func(c echo.Context) error {
		var query QueryParam
		if err := c.Bind(&query); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.String(http.StatusOK, "Not Implemented")
	}
}
