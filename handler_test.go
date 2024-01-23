package yujin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rank1zen/yujin"
	"github.com/rank1zen/yujin/postgresql"
	"github.com/rank1zen/yujin/postgresql/db"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetSummonerRecordsByPuuid(t *testing.T) {
	r := gin.New()
	r.Use(yujin.ErrorHandler())

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	conn := postgresql.NewConnection(t)

	err := conn.Ping(ctx)
	if err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	r.GET("/v1/:puuid", yujin.HandleGetSummonerRecordsByPuuid(db.New(conn)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/testing", nil)
	r.ServeHTTP(w, req)

	log.Print(w.Body)
	assert.Equal(t, 200, w.Code)
}
