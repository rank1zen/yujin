package internal

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/KnutZuidema/golio"
	"github.com/labstack/echo/v4"
	"github.com/rank1zen/yujin/internal/postgres"
)

func GetSummoner(q *postgres.Query) echo.HandlerFunc {
	return func(c echo.Context) error {
		record, err := q.SummonerV4.SelectSummonerRecordNewestByPuuid(context.Background(), c.Param("puuid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, record)
	}
}

func GetSummonerByPuuidRecent(q *postgres.Query) echo.HandlerFunc {
	return func(c echo.Context) error {
		record, err := q.SummonerV4.SelectSummonerRecordNewestByPuuid(context.Background(), c.Param("puuid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, record)
	}
}

func GetSummonerByPuuid(q *postgres.Query) echo.HandlerFunc {
	return func(c echo.Context) error {
		records, err := q.SummonerV4.SelectSummonerRecordsByPuuid(context.Background(), c.Param("puuid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, *records)
	}
}

func GetSummonerByPuuidCount(q *postgres.Query) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		count, err := q.SummonerV4.CountSummonerRecordsByPuuid(ctx, c.Param("puuid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, map[string]int64{"count": count})
	}
}

func PostSummonerByName(q *postgres.Query, gc *golio.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ts := time.Now()

		summoner, err := gc.Riot.LoL.Summoner.GetByName(c.Param("name"))
		if err != nil {
			err = fmt.Errorf("riot error: %w", err)
			return echo.NewHTTPError(http.StatusServiceUnavailable, err.Error())
		}

		record := postgres.SummonerRecordArg{
			Puuid:         summoner.PUUID,
			AccountId:     summoner.AccountID,
			SummonerId:    summoner.ID,
			Name:          summoner.Name,
			ProfileIconId: summoner.ProfileIconID,
			SummonerLevel: summoner.SummonerLevel,
			RevisionDate:  summoner.RevisionDate,
		}

		id, err := q.SummonerV4.InsertSummonerRecord(context.Background(), &record, ts)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusCreated, id)
	}
}

func DeleteSummoner(q *postgres.Query) echo.HandlerFunc {
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

func GetSummonerByNameRecent(q *postgres.Query) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}

func GetSummonerByNameCount(q *postgres.Query) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}

func GetSummonerByName(q *postgres.Query) echo.HandlerFunc {
	return func(c echo.Context) error {
		records, err := q.SummonerV4.SelectSummonerRecordsByName(context.Background(), c.Param("name"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, *records)
	}
}
