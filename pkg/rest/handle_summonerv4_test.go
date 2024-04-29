package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/KnutZuidema/golio"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	var gc *golio.Client
	var q *postgres.Query
	handler := internal.PostSummonerByName(q, gc)

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

func TestGetSummonerRecord(t *testing.T) {
        req, err := http.NewRequest(http.MethodGet, "/summonerv4/", nil)
        require.NoError(t, err)

        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(rest.)
        handler.ServeHTTP(rr, req)

}
