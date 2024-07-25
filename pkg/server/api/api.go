package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rank1zen/yujin/pkg/database"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Router(db *database.DB, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()

	logger = logger.With(zap.String("api version", "0.0.1"))

	router.Use(middleware.NoCache)
	router.Use(loggingMiddleware(logger))

	router.Get("/profile/{puuid}", profileMatchlist(db))
	// router.Get("/profile/{puuid}/{matchID}", profileMatchSummary(db))

	return router
}

func loggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			duration := time.Since(start).Microseconds()

			fields := []zapcore.Field{
				zap.Int64("duration", duration),
				zap.String("method", r.Method),
			}

			logger.Info("", fields...)
		})
	}
}
