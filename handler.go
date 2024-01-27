package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func HandleGetSummonerRecordsByPuuid() echo.HandlerFunc {
	type PathParam struct {
		Puuid string `param:"puuid"`
	}
	return func(c echo.Context) error {
		var path PathParam
		if err := c.Bind(&path); err != nil {
			return echo.NewBindingError("puuid", nil, nil, err)
		}

		return c.String(http.StatusOK, "Not Implemented")
	}
}

func HandleGetSummonerRecordsByName() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Not Implemented")
	}
}

type PostSummonerRecordRequest struct {
	RecordDate    time.Time `json:"record_date" validate:"required"`
	AccountId     string    `json:"account_id"`
	ProfileIconId int32     `json:"profile_icon_id"`
	RevisionDate  int64     `json:"revision_date"`
	Name          string    `json:"name"`
	Puuid         string    `json:"puuid"`
	SummonerLevel int64     `json:"summoner_level"`
}

func HandlePostSummonerRecord() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		var request PostSummonerRecordRequest
		if err = c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err = c.Validate(request); err != nil {
			return err
		}

		return c.String(http.StatusOK, "Not Implemented")
	}
}

func HandleGetSoloqRecordByName() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Not Implemented")
	}
}

func HandlePostSoloqRecord() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Not Implemented")
	}
}

func HandleGetSummonerRecordCountByPuuid() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Not Implemented")
	}
}

func HandleGetSoloqRecordById() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Not Implemented")
	}
}

func HandleGetMatch() echo.HandlerFunc {
	type PathParam struct {
		MatchId string `param:"id"`
	}
	return func(c echo.Context) (err error) {
		var path PathParam
		if err := c.Bind(&path); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		// database

		return c.String(http.StatusOK, "Not Implemented")
	}
}

type PostMatchRequest struct {
	MatchId string `json:"match_id" validate:"required"`
}

func HandlePostMatch() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		var request PostMatchRequest
		if err = c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		// database

		return c.String(http.StatusCreated, "Not Implemented")
	}
}

func HandleHome() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		return c.String(http.StatusOK, "Welcome to YUJIN.GG")
	}
}
