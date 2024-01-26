package main

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, p *pgxpool.Pool) {
	g := e.Group("/v1/")
	g.GET("healthcheck", HandleHealth(p))
}
