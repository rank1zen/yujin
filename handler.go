package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rank1zen/yujin/postgresql"
)

type SummonerProfileQuery struct {
	Name string
}

type SummonerProfileBody struct {
	Name       string `json:"name" validate:"required"`
	Puuid      string `json:"puuid" validate:"required"`
	AccountId  string `json:"account_id" validate:"required"`
	SummonerId string `json:"summoner_id" validate:"required"`
}

func HandleHome() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		return c.JSON(http.StatusOK, "Welcome to YUJIN.GG")
	}
}

func HandleGetSummonerProfile(q *postgresql.Queries) echo.HandlerFunc {
	type PathParam struct {
		Name string `param:"name"`
	}

	return func(c echo.Context) error {
		var path PathParam
		if err := c.Bind(&path); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, "Not Implemented")
	}
}

func HandlePostSummonerProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		var body SummonerProfileBody
		err := c.Bind(&body)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, body)
	}
}
