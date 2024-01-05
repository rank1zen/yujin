package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rank1zen/yujin/internal/postgresql/db"
	"github.com/rank1zen/yujin/internal/riotgrpc/proto"
)

type summonerHandler struct {
	q *db.Queries
	rc proto.RiotSummonerClient
}

func (s *summonerHandler) getAll(c *gin.Context) {
	args := db.SelectSummonerRecordsParams{
		Limit: 10,
		Offset: 0,
	}

	r, err := q.SelectSummonerRecords(c, args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err})
		return
	}

	c.JSON(http.StatusOK, r)
}
