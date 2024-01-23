package yujin

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"github.com/rank1zen/yujin/postgresql/db"
)

type Config struct {
	Debug              bool   `required:"true"`
	PostgresConnString string `required:"true" split_words:"true"`
	Addr               string `required:"true"`
}

func main() {
	var conf Config
	err := envconfig.Process("YUJIN", &conf)
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.New(context.Background(), conf.PostgresConnString)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.Use(ErrorHandler())
	v1 := router.Group("/v1")

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Welcome YUJIN.GG")
	})

	router.GET("/ready", func(c *gin.Context) {
		err := pool.Ping(c)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": err.Error()})
			return
		}
		c.JSON(http.StatusOK, "db is up")
	})

	v1.GET("/get", HandleGetSummonerRecordsByName(db.New(pool)))
	v1.POST("/renew", HandlePostSummonerRecord(db.New(pool)))

	router.Run()

}
