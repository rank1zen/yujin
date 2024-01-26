package main_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/rank1zen/yujin"
)

func TestPostMatch(t *testing.T) {
	tests := []struct {
		payload string
		httpStatusCode int
	}{
		{payload: `{"match_id":1}`, httpStatusCode: http.StatusServiceUnavailable},
		{payload: `{"match_id":"TESTINGTESTING"}`, httpStatusCode: http.StatusCreated},
	}

	e := echo.New()
	h := main.HandlePostMatch()

	for _, tc := range tests {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.payload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, h(c)) {
			assert.Equal(t, tc.httpStatusCode, rec.Code)
		}
	}
}
