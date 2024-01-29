package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rank1zen/yujin/postgresql"
)

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

func HandlePostSummonerProfile(q *postgresql.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		var body SummonerProfileBody
		if err := c.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusOK, "Not Implemented")
	}
}
