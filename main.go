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
	"github.com/rank1zen/yujin/postgresql"
	"go.uber.org/zap"
)

func main() {
	log := zap.Must(zap.NewProduction())
	defer log.Sync()

	e := echo.New()

	logcfg := middleware.RequestLoggerConfig{
		LogHost:     true,
		LogURI:      true,
		LogStatus:   true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			fields := []zap.Field{
				zap.String("host", v.Host),
				zap.String("uri", v.URI),
				zap.Int("status", v.Status),
				zap.Error(v.Error),
			}

			switch {
			case v.Status >= 500:
				log.Error("server error", fields...)
			case v.Status >= 400:
				log.Warn("client error", fields...)
			case v.Status >= 300:
				log.Info("redirection", fields...)
			default:
				log.Info("success", fields...)
			}

			return nil
		},
	}

	e.Use(middleware.RequestLoggerWithConfig(logcfg))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		ErrorMessage: "request exceeded timeout",
		Timeout:      10 * time.Second,
	}))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := postgresql.BackoffRetryPool(ctx, os.Getenv("YUJIN_PG_STRING"))
	if err != nil {
		log.Warn("can't connect to database")
	}

	e.GET("/", HandleHome())
	RegisterRoutes(e, pool)

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
