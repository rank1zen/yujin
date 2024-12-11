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
	"github.com/rank1zen/yujin/internal"
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

type contextKey string

const requestID contextKey = "request_id"

func subPartials(db *database.DB) *chi.Mux {
	mux := chi.NewMux()

	middlewareChain := []func(http.Handler) http.Handler{
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.Header.Get("HX-Request")
				w.Header().Get("HX-Request")
			})
		},
	}

	mux.Use(middlewareChain...)

	mux.Get("/profile/{puuid}/matchlist", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		puuid := chi.URLParam(r, "puuid")
		page := request.QueryIntParam(r, "page", 0)

		m, err := db.ProfileGetMatchList(ctx, riot.PUUID(puuid), page, true)
		switch err {
		case nil:
			html.OK(w, r, partials.ProfileMatchList(m))
		default:
			html.ServerError(w, r, partials.ProfileMatchListError(), err)
		}
	})

	mux.Get("/profile/{puuid}/champstats", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		puuid := chi.URLParam(r, "puuid")

		m, err := db.ProfileGetChampionStatList(ctx, riot.PUUID(puuid), internal.Season2020)

		switch {
		case err == nil:
			html.OK(w, r, partials.ProfileChampionStatList(m))
		case errors.As(err, 1):
			html.OK(w, r, partials.ProfileLiveGameNotFoundError())
		default:
			html.ServerError(w, r, partials.ProfileLiveGameError(), err)
		}
	})

	mux.Get("/profile/{puuid}/live", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		puuid := chi.URLParam(r, "puuid")
		m, err := db.ProfileGetLiveGame(ctx, puuid)

		switch {
		case err == nil:
			html.OK(w, r, partials.ProfileLiveGame(m))
		case errors.As(err, 1):
			html.OK(w, r, partials.ProfileLiveGameNotFoundError())
		default:
			html.ServerError(w, r, partials.ProfileLiveGameError(), err)
		}
	})

	return mux
}

func subPages(db *database.DB) *chi.Mux {
	mux := chi.NewRouter()

	middlewareChain := []func(http.Handler) http.Handler{
		middleware.NoCache,
	}

	mux.Use(middlewareChain...)

	mux.Get("/profile/{puuid}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		puuid := chi.URLParam(r, "puuid")

		exists, err := db.ProfileExists(ctx, puuid)
		if err != nil {
			html.ServerError(w, r, pages.InternalServerError(), err)
			return
		}

		if !exists {
			html.BadRequest(w, r, pages.ProfileDoesNotExist(), err)
			return
		}

		profile, err := db.ProfileGetHeader(ctx, puuid)
		if err != nil {
			html.ServerError(w, r, pages.InternalServerError(), err)
			return
		}

		html.OK(w, r, pages.Profile(profile))
	})

	mux.Post("/profile/{puuid}/update", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		puuid := chi.URLParam(r, "puuid")

		err := db.ProfileUpdate(ctx, puuid)
		if err != nil {
			html.ServerError(w, r, nil, fmt.Errorf("updating summary: %w", err))
			return
		}

		w.Header().Set("HX-Refresh", "true")
		w.WriteHeader(http.StatusOK)
	})

	return mux
}

func Routes(db *database.DB) *chi.Mux {
	base := chi.NewRouter()

	middlewareChain := []func(http.Handler) http.Handler{
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

	base.Use(middlewareChain...)

	base.NotFound(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger := logging.FromContext(ctx)

		logger.Warn(r.URL.String())
		templ.Handler(pages.NotFound(), templ.WithStatus(http.StatusNotFound)).ServeHTTP(w, r)
	})

	base.Mount("/", subPages(db))

	base.Mount("/partials", subPartials(db))

	return base
}
