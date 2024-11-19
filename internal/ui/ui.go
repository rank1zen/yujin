package ui

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/rank1zen/yujin/internal/database"
	"github.com/rank1zen/yujin/internal/http/request"
	"github.com/rank1zen/yujin/internal/http/response/html"
	"github.com/rank1zen/yujin/internal/logging"
	"github.com/rank1zen/yujin/internal/riot"
	"github.com/rank1zen/yujin/internal/ui/pages"
	"github.com/rank1zen/yujin/internal/ui/partials"
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
			"/profile/{puuid}",
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				puuid := chi.URLParam(r, "puuid")
				exists, err := db.ProfileExists(ctx, puuid)
				if err != nil {
					html.ServerError(w, r, pages.ProfileNotFound(puuid), err)
					// internal server err html
					return
				}

				if exists {
					profile, err := db.ProfileGetHeader(ctx, puuid)
					if err != nil {
						html.ServerError(w, r, pages.ProfileNotFound(puuid), err)
						return
					}

					html.OK(w, r, pages.Profile(profile, puuid))
				} else {
					html.BadRequest(w, r, pages.ProfileNotFound(puuid), err)
				}
			},
		},
		{
			"/profile/{name}/matchlist",
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				name := chi.URLParam(r, "name")
				page := request.QueryIntParam(r, "page", 0)
				m, err := db.ProfileGetMatchList(ctx, name, page, true)
				switch err {
				case nil:
					html.OK(w, r, partials.ProfileMatchList(m))
				default:
					html.ServerError(w, r, partials.ProfileMatchListError(), err)
				}
			},
		},
		{
			"/profile/{name}/live",
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				name := chi.URLParam(r, "name")
				m, err := db.ProfileGetLiveGame(ctx, name)
				switch {
				case err == nil:
					html.OK(w, r, partials.ProfileLiveGame(m))
				case errors.As(err, 1):
					html.OK(w, r, partials.ProfileLiveGameNotFoundError())
				default:
					html.ServerError(w, r, partials.ProfileLiveGameError(), err)
				}
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
				err := db.ProfileUpdate(ctx, name)
				if err != nil {
					html.ServerError(w, r, nil, fmt.Errorf("updating summary: %w", err))
					return
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

func GenMatchListQuery(puuid riot.PUUID, page int) string {
	return fmt.Sprintf("/profile/%s/matchlist?page=%d", puuid, page)
}

func GenLiveGameQuery(puuid riot.PUUID) string {
	return fmt.Sprintf("/profile/%s/livegame", puuid)
}

func GenChampionStatsQuery(puuid riot.PUUID) string {
	return fmt.Sprintf("/profile/%s/matchlist", puuid)
}
