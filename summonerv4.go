package main

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rank1zen/yujin/postgresql"
)

type SummonerRecordBody struct {
	RecordDate    time.Time `json:"record_date" validate:"required"`
	Name          string    `json:"name" validate:"required"`
	Puuid         string    `json:"puuid" validate:"required"`
	AccountId     string    `json:"account_id" validate:"required"`
	SummonerId    string    `json:"summoner_id" validate:"required"`
	SummonerLevel int       `json:"summoner_level" validate:"required"`
	ProfileIconId int       `json:"profile_icon_id" validate:"required"`
	RevisionDate  int       `json:"revision_date" validate:"required"`
}

func GetSummoner(q *postgresql.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		record, err := q.SelectSummonerRecord(ctx, c.Param("uuid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, record)
	}
}

func GetSummonerByPuuidRecent(q *postgresql.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		record, err := q.SelectSummonerRecentByPuuid(ctx, c.Param("puuid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusServiceUnavailable, err.Error())
		}

		return c.JSON(http.StatusOK, record)
	}
}

func GetSummonerByPuuid(q *postgresql.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		records, err := q.SelectSummonerRecordsByPuuid(ctx, c.Param("puuid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusServiceUnavailable, err.Error())
		}

		return c.JSON(http.StatusOK, *records)
	}
}

func HandleGetSummonerRecordsByName(q *postgresql.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		records, err := q.SelectSummonerRecordsByName(ctx, c.Param("name"))
		if err != nil {
			return echo.NewHTTPError(http.StatusServiceUnavailable, err.Error())
		}

		return c.JSON(http.StatusOK, *records)
	}
}

func GetSummonerByPuuidCount(q *postgresql.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		count, err := q.CountSummonerRecordsByPuuid(ctx, c.Param("puuid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusServiceUnavailable, err.Error())
		}

		return c.JSON(http.StatusOK, map[string]int64{"count":count})
	}
}

func PostSummoner(q *postgresql.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		var body SummonerRecordBody
		err := c.Bind(&body)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = c.Validate(body)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		id, err := q.InsertSummonerRecord(ctx, &postgresql.SummonerRecordArg{}, body.RecordDate)
		if err != nil {
			return echo.NewHTTPError(http.StatusServiceUnavailable, err.Error())
		}

		return c.JSON(http.StatusCreated, id)
	}
}

func DeleteSummoner(q *postgresql.Queries) echo.HandlerFunc {
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

func GetSummonerByNameRecent(q *postgresql.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}

func GetSummonerByNameCount(q *postgresql.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func GetSummonerByName(q *postgresql.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
