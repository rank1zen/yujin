package main

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rank1zen/yujin/postgresql"
	"go.uber.org/zap"
)

func MiddleZapLogger(l *zap.Logger) echo.MiddlewareFunc {
	conf := middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			fields := []zap.Field{
				zap.Int("status", v.Status),
			}

			switch {
			case v.Status >= 500:
				l.Error("server error", fields...)
			case v.Status >= 400:
				l.Warn("client error", fields...)
			case v.Status >= 300:
				l.Info("redirection", fields...)
			default:
				l.Info("success", fields...)
			}

			return nil
		},
	}

	return middleware.RequestLoggerWithConfig(conf)
}

func MiddleDbHealth(p *pgxpool.Pool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			if err := postgresql.CheckPool(ctx, p); err != nil {
				return echo.NewHTTPError(http.StatusServiceUnavailable, err.Error())
			}
			return next(c)
		}
	}
}
