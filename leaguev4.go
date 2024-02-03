package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)


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

func HandleGetSoloqRecordById() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Not Implemented")
	}
}
