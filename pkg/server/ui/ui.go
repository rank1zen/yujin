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
	"github.com/rank1zen/yujin/pkg/database"
	"github.com/rank1zen/yujin/pkg/http/request"
	"github.com/rank1zen/yujin/pkg/http/response/html"
	"github.com/rank1zen/yujin/pkg/logging"
	"github.com/rank1zen/yujin/pkg/server/ui/pages"
	"github.com/rank1zen/yujin/pkg/server/ui/partials"
	"github.com/rank1zen/yujin/pkg/server/ui/static"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Routes(db *database.DB) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.NoCache)
	router.Use(loggingMiddleware)
	router.Use(middleware.Recoverer)

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger := logging.FromContext(ctx)

		logger.Warn(r.URL.String())
		templ.Handler(pages.NotFound(), templ.WithStatus(http.StatusNotFound)).ServeHTTP(w, r)
	})

	for _, handler := range []struct {
		path string
		fn   http.HandlerFunc
	}{
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
			path: "/",
			fn: func(w http.ResponseWriter, r *http.Request) {
				html.OK(w, r, pages.About())
			},
		},
		{
			path: "/profile/{name}",
			fn: func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()

				nameParam := chi.URLParam(r, "name")
				name, err := database.ParseRiotName(nameParam)
				if err != nil {
					html.BadRequest(w, r, pages.NotFound(), err)
				}

				profile, err := db.GetProfileSummary(ctx, name)
				if err != nil {
					html.ServerError(w, r, pages.ProfileNotFound(name), err)
				}

				html.OK(w, r, pages.Profile(r, profile))
			},
		},
		{
			"/profile/{name}/update",
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()

				nameParam := chi.URLParam(r, "name")
				name, err := database.ParseRiotName(nameParam)
				if err != nil {
					html.BadRequest(w, r, pages.NotFound(), err)
				}

				profile, err := db.UpdateProfileSummary(ctx, name)
				if err != nil {
					html.ServerError(w, r, nil, fmt.Errorf("updating summary: %w", err))
				}

				matchlist, err := db.UpdateProfileMatchlist(ctx, name, 0)
				if err != nil {
					html.ServerError(w, r, nil)
				}
			},
		},
		{
			path: "/profile/{name}/matchlist",
			fn: func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()

				name, err := database.ParseRiotName(chi.URLParam(r, "name"))
				if err != nil {
					html.BadRequest(w, r, partials.ProfileMatchSummaryError(), err)
					return
				}

				page := request.QueryIntParam(r, "page", 10)
				matches, err := db.GetProfileMatchList(ctx, name, page)
				if err != nil {
					html.ServerError(w, r, partials.ProfileMatchlistError(), err)
					return
				}

				html.OK(w, r, partials.ProfileMatchlist(r, matches))
			},
		},
		{
			path: "/profile/{name}/matchlist/{matchid}",
			fn: func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()

				name, err := database.ParseRiotName(chi.URLParam(r, "name"))
				if err != nil {
					html.BadRequest(w, r, partials.ProfileMatchSummaryError(), err)
					return
				}

				matchId := database.RiotMatchId(chi.URLParam(r, "matchid"))
				match, err := db.GetProfileMatchSummary(ctx, name, matchId)
				if err != nil {
					html.ServerError(w, r, partials.ProfileMatchSummary(match), err)
					return
				}

				html.OK(w, r, partials.ProfileMatchSummary(match))
			},
		},
	} {
		router.Get(handler.path, handler.fn)
	}

	return router
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
