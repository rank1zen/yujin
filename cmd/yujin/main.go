package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/postgresql"
	"go.uber.org/zap"
)

func main() {
	log := zap.Must(zap.NewDevelopment())
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	pool, err := postgresql.NewConnectionPool(context.Background())

	gc := golio.NewClient(os.Getenv("YUJIN_RIOT_API_KEY"), golio.WithRegion(api.RegionNorthAmerica))

	internal.Routes(e, pool, gc)


	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		err := e.Start(fmt.Sprintf("0.0.0.0:%s", os.Getenv("YUJIN_PORT")))
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("shutting down server...")
		}
	}()

	<-ctx.Done()

	err = e.Shutdown(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
}
