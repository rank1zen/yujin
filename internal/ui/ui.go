package ui

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/rank1zen/yujin/internal/database"
	"github.com/rank1zen/yujin/internal/http/request"
	"github.com/rank1zen/yujin/internal/http/response/html"
	"github.com/rank1zen/yujin/internal/logging"
	"github.com/rank1zen/yujin/internal/ui/pages"
	"github.com/rank1zen/yujin/internal/ui/partials"
	"github.com/rank1zen/yujin/internal/ui/static"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type pathHandler struct {
	path string
	fn   http.HandlerFunc
}

type contextKey string

const requestID contextKey = "request_id"

func Routes(db *database.DB) *chi.Mux {
	router := chi.NewRouter()

	middlewareChain := []func(http.Handler) http.Handler{
		middleware.NoCache,
		middleware.Recoverer,
		func(next http.Handler) http.Handler {
			fn := func(w http.ResponseWriter, r *http.Request) {
				logger := logging.Get()

				id := uuid.Must(uuid.NewRandom()).String()

				ctx := context.WithValue(r.Context(), requestID, id)
				r = r.WithContext(ctx)

				logger = logger.With(zap.String("request_id", id))
				r = r.WithContext(logging.WithContext(ctx, logger))

				w.Header().Set("X-Yujin-Request-ID", id)
				wrw := middleware.NewWrapResponseWriter(w, 1)

				defer func(start time.Time) {
					fields := []zapcore.Field{
						zap.Duration("duration_ms", time.Since(start)),
						zap.Int("response_bytes", wrw.BytesWritten()),
						zap.Int("status", wrw.Status()),
						zap.String("url", r.RequestURI),
						zap.String("user_agent", r.UserAgent()),
					}

					switch wrw.Status() {
					case http.StatusOK:
						logger.Info(r.Method, fields...)
					default:
						logger.Error(r.Method, fields...)
					}
				}(time.Now())

				next.ServeHTTP(wrw, r)
			}

			return http.HandlerFunc(fn)
		},
	}

	router.Use(middlewareChain...)

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger := logging.FromContext(ctx)

		logger.Warn(r.URL.String())
		templ.Handler(pages.NotFound(), templ.WithStatus(http.StatusNotFound)).ServeHTTP(w, r)
	})

	for _, handler := range []pathHandler{
		{
			"/static/*",
			func(w http.ResponseWriter, r *http.Request) {
				rctx := chi.RouteContext(r.Context())
				pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
				fs := http.StripPrefix(pathPrefix, http.FileServer(http.FS(static.StylesheetFiles)))
				fs.ServeHTTP(w, r)
			},
		},
		{
			"/",
			func(w http.ResponseWriter, r *http.Request) {
				html.OK(w, r, pages.About())
			},
		},
		{
			"/profile/{name}",
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				name := chi.URLParam(r, "name")
				profile, err := db.GetProfileSummary(ctx, name)
				if err != nil {
					html.ServerError(w, r, pages.ProfileNotFound(r), err)
				}

				html.OK(w, r, pages.ProfileMatchList(r, profile))
			},
		},
		{
			"/profile/{name}/other",
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				name := chi.URLParam(r, "name")
				profile, err := db.GetProfileSummary(ctx, name)
				if err != nil {
					html.ServerError(w, r, pages.ProfileNotFound(r), err)
				}

				html.OK(w, r, pages.ProfileOther(r, profile))
			},
		},
		{
			"/profile/{name}/matchlist",
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				name := chi.URLParam(r, "name")
				page := request.QueryIntParam(r, "page", 0)
				matches, err := db.GetProfileMatchList(ctx, name, page, true)
				if err != nil {
					html.ServerError(w, r, partials.ProfileMatchListError(), err)
					return
				}

				html.OK(w, r, partials.ProfileMatchList(r, matches))
			},
		},
	} {
		router.Get(handler.path, handler.fn)
	}

	for _, handler := range []pathHandler{
		{
			"/profile/{name}/update",
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				name := chi.URLParam(r, "name")
				err := db.UpdateProfile(ctx, name)
				if err != nil {
					html.ServerError(w, r, nil, fmt.Errorf("updating summary: %w", err))
				}

				w.Header().Set("HX-Refresh", "true")
				w.WriteHeader(http.StatusOK)
			},
		},
	} {
		router.Post(handler.path, handler.fn)
	}

	return router
}
