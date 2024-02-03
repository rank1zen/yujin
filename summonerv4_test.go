package main_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rank1zen/yujin"
	"github.com/rank1zen/yujin/postgresql"
	"github.com/stretchr/testify/assert"
)

func TestGetSummonerRecordsByPuuid(t *testing.T) {
}

func TestPostSummonerRecord(t *testing.T) {
	tests := []struct {
		payload        string
		httpStatusCode int
	}{
		{payload: `{}`, httpStatusCode: http.StatusBadRequest},
		{payload: `{"name":"JOJO"}`, httpStatusCode: http.StatusBadRequest},
		{
			payload: `{"record_date":"","puuid":"","account_id":"","profile_icon_id":"","name":"","summoner_level":"","revision_date":""}`,
			httpStatusCode: http.StatusBadRequest,
		},
	}


	e := echo.New()
	e.Validator = main.NewValidator()
	handler := main.PostSummoner(q)

	for _, tc := range tests {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.payload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, handler(c)) {
			assert.Equal(t, tc.httpStatusCode, rec.Code)
		}
	}
}

func handlePostSummonerRecord() echo.HandlerFunc {
	return func(c echo.Context) error {
		var body main.SummonerRecordBody
		if err := c.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := c.Validate(body); err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, id)
	}
}
