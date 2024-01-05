package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/yujin/internal/postgresql"
	"github.com/rank1zen/yujin/internal/postgresql/db"
	"github.com/rank1zen/yujin/internal/riotgrpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var q *db.Queries

var rc proto.RiotSummonerClient

func main() {
	r := gin.Default()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	dp, err := pgxpool.New(ctx, "")
	if err != nil {
		panic(err)
	}

	ccOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	cc, err := grpc.DialContext(ctx, "", ccOpts...)
	if err != nil {
		panic(err)
	}

	q = db.New(dp)
	rc = proto.NewRiotSummonerClient(cc)


	v1 := r.Group("/api/v1/summoner")
	{
		v1.GET("/get/all", getAll)
		v1.GET("/get/all/newest", notImplemented)
		v1.GET("/get/by-puuid/:puuid", notImplemented)
		v1.GET("/get/by-puuid/:puuid/newest", notImplemented)
		v1.POST("/renew/by-name/:name", renewByName)
	}

	r.Run()
}

func getAll(c *gin.Context) {
	offset := c.DefaultQuery("offset", "0")
	limit := c.DefaultQuery("limit", "10")

	off, err := strconv.Atoi(offset)
	if err != nil {
		off = 0
	}

	lim, err := strconv.Atoi(limit)
	if err != nil {
		lim = 10
	}

	args := db.SelectSummonerRecordsParams{
		Limit: int32(lim),
		Offset: int32(off),
	}

	r, err := q.SelectSummonerRecords(c, args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err})
		return
	}

	c.JSON(http.StatusOK, r)
}

func renewByName(c *gin.Context) {
	name := c.Param("name")

	rpcreq := proto.ByNameRequest{Name: name}

	s, err := rc.ByName(c, &rpcreq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	args := db.InsertSummonerParams{
		Puuid: s.GetPuuid(),
		AccountID: s.GetAccountId(),
		ID: s.GetSummonerId(),
		SummonerLevel: s.GetLevel(),
		ProfileIconID: s.GetProfileIconId(),
		Name: s.GetName(),
		RevisionDate: s.GetLastRevision(),
		RecordDate: postgresql.NewTimestamp(time.Now()),
	}

	uuid, err := q.InsertSummoner(c, args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, postgresql.UUIDString(uuid))
}

func notImplemented(c *gin.Context) {
	c.JSON(http.StatusOK, "not implemented")
}
