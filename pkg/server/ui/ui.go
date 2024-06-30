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

type ui struct {
	db     *database.DB
	logger *zap.Logger
}

func NewUI(db *database.DB, logger *zap.Logger) *ui {
	return &ui{
		db:     db,
		logger: logger,
	}
}

func (ui *ui) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.NoCache)
	router.Use(ui.loggingMiddleware(ui.logger))
	router.Use(ui.requestIdMiddleware)

	// STATIC FILES
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	FileServer(router, "/static", filesDir)

	router.Get("/", ui.home())
	router.Get("/profile/{puuid}", ui.profile())

	return router
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

func (a *ui) requestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.Must(uuid.NewRandom()).String()
		w.Header().Set("X-Yujin-Request-Id", id)
		next.ServeHTTP(w, r)
	})
}

func (ui *ui) loggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			duration := time.Since(start).Milliseconds()

			fields := []zapcore.Field{
				zap.Int64("duration", duration),
				zap.String("method", r.Method),
			}

			logger.Info("", fields...)
		})
	}
}
