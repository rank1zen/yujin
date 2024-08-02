package ui

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/rank1zen/yujin/pkg/database"
	"github.com/rank1zen/yujin/pkg/http/request"
	"github.com/rank1zen/yujin/pkg/logging"
	"github.com/rank1zen/yujin/pkg/server/ui/pages"
	"github.com/rank1zen/yujin/pkg/server/ui/partials"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Routes(db *database.DB, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.NoCache)
	router.Use(loggingMiddleware)
	router.Use(middleware.Recoverer)

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(pages.NotFound(), templ.WithStatus(http.StatusNotFound)).ServeHTTP(w, r)
	})

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	FileServer(router, "/static", filesDir)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(pages.About(), templ.WithStatus(200)).ServeHTTP(w, r)
	})

	router.Get("/profile/{name}", func(w http.ResponseWriter, r *http.Request) {
		fn := func() *templ.ComponentHandler {
			ctx := r.Context()

			logger := logging.FromContext(ctx)
			nameParam := chi.URLParam(r, "name")

			name, err := database.ParseRiotName(nameParam)
			if err != nil {
				logger.Warn(http.StatusText(http.StatusBadRequest), zap.Error(err))
				return templ.Handler(pages.NotFound(), templ.WithStatus(http.StatusBadRequest))
			}

			profile, err := db.UpdateProfileSummary(ctx, name)
			if err != nil {
				logger.Warn(http.StatusText(http.StatusBadRequest), zap.Error(err))
				return templ.Handler(pages.ProfileNotFound(name), templ.WithStatus(http.StatusInternalServerError))
			}

			return templ.Handler(pages.Profile(r, profile), templ.WithStatus(http.StatusOK))
		}

		fn().ServeHTTP(w, r)
	})

	router.Get("/profile/{name}/matchlist", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		name, err := database.ParseRiotName(chi.URLParam(r, "name"))
		if err != nil {
			templ.Handler(partials.ProfileMatchlistError(), templ.WithStatus(http.StatusBadRequest)).ServeHTTP(w, r)
			return
		}

		page := request.QueryIntParam(r, "page", 10)

		matches, err := db.GetProfileMatchList(ctx, name, page)
		if err != nil {
			templ.Handler(partials.ProfileMatchlistError(), templ.WithStatus(http.StatusInternalServerError)).ServeHTTP(w, r)
			return
		}

		templ.Handler(partials.ProfileMatchlist(matches), templ.WithStatus(http.StatusOK)).ServeHTTP(w, r)
	})

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

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logging.Get()

		id := uuid.Must(uuid.NewRandom()).String()

		ctx := context.WithValue(r.Context(), "request_id", id)
		r = r.WithContext(ctx)

		logger.With(zap.String("request_id", id))

		w.Header().Set("X-Yujin-Request-ID", id)

		wrw := middleware.NewWrapResponseWriter(w, 1)

		r = r.WithContext(logging.WithContext(ctx, logger))

		defer func(start time.Time) {
			fields := []zapcore.Field{
				zap.Duration("duration_ms", time.Since(start)),
				zap.Int("response_bytes", wrw.BytesWritten()),
				zap.Int("status", wrw.Status()),
				zap.String("method", r.Method),
				zap.String("url", r.RequestURI),
				zap.String("user_agent", r.UserAgent()),
			}

			logger.Info("REQUEST", fields...)
		}(time.Now())

		next.ServeHTTP(wrw, r)
	})
}
