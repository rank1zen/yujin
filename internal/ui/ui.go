package ui

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rank1zen/yujin/internal/database"
	"github.com/rank1zen/yujin/internal/logging"
	"github.com/rank1zen/yujin/internal/ui/pages"
)

type contextKey string

const requestID contextKey = "request_id"

func pagesSubgroup(db *database.DB) *chi.Mux {
	mux := chi.NewRouter()

	middlewareChain := []func(http.Handler) http.Handler{
		middleware.NoCache,
	}

	mux.Use(middlewareChain...)

	mux.Get("/profile/{puuid}", profileShow(db))
	mux.Post("/profile/{puuid}/update", profileRefresh(db))

	return mux
}

func partialsSubroute(db *database.DB) *chi.Mux {
	mux := chi.NewMux()

	mux.Use(checkHTMX)

	mux.Get("/profile/{puuid}/matchlist", profileShowMatchList(db))
	mux.Get("/profile/{puuid}/champstats", profileShowChampStats(db))
	mux.Get("/profile/{puuid}/live", profileShowLiveGame(db))

	return mux
}

func Routes(db *database.DB) *chi.Mux {
	base := chi.NewRouter()

	base.Use(addRequestID, addLoggedFields)

	base.NotFound(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger := logging.FromContext(ctx)

		logger.Warn(r.URL.String())
		templ.Handler(pages.NotFound(), templ.WithStatus(http.StatusNotFound)).ServeHTTP(w, r)
	})

	base.Mount("/", pagesSubgroup(db))
	base.Mount("/partials", partialsSubroute(db))

	return base
}
