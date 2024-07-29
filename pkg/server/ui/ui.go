package ui

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/rank1zen/yujin/pkg/database"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Routes(db *database.DB, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.NoCache)
	router.Use(requestIdMiddleware)
	router.Use(loggingMiddleware(logger))
	router.Use(middleware.Recoverer)

	router.NotFound(NotFoundHandler())

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	FileServer(router, "/static", filesDir)

	router.Get("/", aboutPage)

	router.Get("/profile/{name}", profilePage(db))
	router.Get("/profile/{name}/matchlist", profileMatchList(db))

	return router
}

func NotFoundHandler() http.HandlerFunc {
	fn := ErrorNotFound()
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func loggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrw := middleware.NewWrapResponseWriter(w, 1)

			next.ServeHTTP(wrw, r)

			duration := time.Since(start).Milliseconds()

			fields := []zapcore.Field{
				zap.Int64("duration_ms", duration),
				zap.String("method", r.Method),
				zap.Int("response#bytes", wrw.BytesWritten()),
				zap.Int("status", wrw.Status()),
				zap.String("url", r.RequestURI),
				zap.String("request#id", wrw.Header().Get("X-Yujin-Request-Id")),
			}

			if wrw.Status() == 200 {
				logger.Info("", fields...)
			} else {
				err := wrw.Header().Get("X-Yujin-Error")
				logger.Error(err, fields...)
			}
		})
	}
}

func requestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.Must(uuid.NewRandom()).String()
		w.Header().Set("X-Yujin-Request-Id", id)
		next.ServeHTTP(w, r)
	})
}
