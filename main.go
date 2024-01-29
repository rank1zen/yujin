package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rank1zen/yujin/postgresql"
	"go.uber.org/zap"
)

func main() {
	log := zap.Must(zap.NewDevelopment())
	defer log.Sync()

	conf, err := LoadConfig()
	if err != nil {
		log.Warn("can't load config. defaulting to preset")
	}

	e := echo.New()
	e.Use(MiddleZapLogger(log))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := postgresql.BackoffRetryPool(ctx, conf.PostgresConnString, log)
	if err != nil {
		log.Warn("can't connect to database")
	}

	e.GET("/", HandleHome())
	RegisterRoutes(e, pool)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		err := e.Start(fmt.Sprintf(":%d", conf.ServerPort))
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("shutting down server...")
		}
	}()

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err.Error())
	}
}
