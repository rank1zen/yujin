package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/rank1zen/yujin/postgresql"
)

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)

	conf, err := LoadConfig()
	if err != nil {
		e.Logger.Fatal("can't load config")
	}

	e.Use(middleware.Logger())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	e.Logger.Info("trying to connect to postgresql")
	pool, err := postgresql.NewPool(ctx, conf.PostgresConnString)
	if err != nil {
		e.Logger.Fatal("can't make a pool")
	}

	e.GET("/", HandleHome())
	RegisterRoutes(e, pool)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		err := e.Start(fmt.Sprintf(":%d", conf.ServerPort))
		if err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
