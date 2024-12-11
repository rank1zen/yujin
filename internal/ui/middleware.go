package ui

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/rank1zen/yujin/internal/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func addRequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		id := uuid.Must(uuid.NewRandom()).String()
		ctx := context.WithValue(r.Context(), requestID, id)
		r = r.WithContext(ctx)

		w.Header().Set("X-Yujin-Request-ID", id)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func addLoggedFields(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		wrw := middleware.NewWrapResponseWriter(w, 1)

		start := time.Now()

		next.ServeHTTP(wrw, r)

		fields := []zapcore.Field{
			zap.Duration("duration_ms", time.Since(start)),
			zap.Int("response_bytes", wrw.BytesWritten()),
			zap.Int("status", wrw.Status()),
			zap.String("url", r.RequestURI),
			zap.String("user_agent", r.UserAgent()),
		}

		logger := logging.Get()
		switch wrw.Status() {
		case http.StatusOK:
			logger.Info(r.Method, fields...)
		default:
			logger.Error(r.Method, fields...)
		}
	}

	return http.HandlerFunc(fn)
}
