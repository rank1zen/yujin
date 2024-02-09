package internal

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type MatchBody struct {
	MatchId string `json:"match_id" validate:"required"`
}

func HandleGetMatch() echo.HandlerFunc {
	type PathParam struct {
		MatchId string `param:"match_id"`
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

func HandlePostMatch() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		var request MatchBody
		if err = c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		// database

		return c.String(http.StatusCreated, "Not Implemented")
	}
}
