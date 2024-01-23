package yujin

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rank1zen/yujin/postgresql"
	"github.com/rank1zen/yujin/postgresql/db"
)

func HandleGetSummonerRecordsByPuuid(q *db.Queries) gin.HandlerFunc {
	type Uri struct {
		Puuid string `uri:"puuid" binding:"required"`
	}
	return func(c *gin.Context) {
		var uri Uri
		err := c.ShouldBindUri(&uri)
		if err != nil {
			c.Error(err)
			return
		}

		qctx, cancel := context.WithTimeout(c, 1*time.Second)
		defer cancel()

		params := db.SelectSummonerRecordsByPuuidParams{
			Puuid:  pgtype.Text{String: uri.Puuid, Valid: true},
			Limit:  20,
			Offset: 0,
		}

		r, err := q.SelectSummonerRecordsByPuuid(qctx, params)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, r)
	}
}

func HandleGetSummonerRecordsByName(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		qctx, cancel := context.WithTimeout(c, 1*time.Second)
		defer cancel()

		_ = c.Param("name")
		params := db.SelectSummonerRecordsByNameParams{
			Limit:  20,
			Offset: 0,
		}

		r, err := q.SelectSummonerRecordsByName(qctx, params)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, r)
	}
}

func HandlePostSummonerRecord(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params db.InsertSummonerRecordParams
		err := c.ShouldBindJSON(&params)
		if err != nil {
			c.Error(err)
			return
		}

		qctx, cancel := context.WithTimeout(c, 1*time.Second)
		defer cancel()

		uuid, err := q.InsertSummonerRecord(qctx, params)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, postgresql.UUIDString(uuid))
	}
}

func HandleGetSoloqRecordByName(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func HandlePostSoloqRecord(q *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		qctx, cancel := context.WithTimeout(c, 1*time.Second)
		defer cancel()

		var params db.InsertSoloqRecordParams
		err := c.ShouldBindJSON(&params)
		if err != nil {
			c.Error(err)
			return
		}

		uuid, err := q.InsertSoloqRecord(qctx, params)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, postgresql.UUIDString(uuid))
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, err := range c.Errors {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"pgx error": err})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"unknown error": err})
			}
		}
	}
}
