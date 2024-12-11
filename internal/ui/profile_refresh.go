package ui

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rank1zen/yujin/internal/database"
	"github.com/rank1zen/yujin/internal/http/response/html"
)

func profileRefresh(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		puuid := chi.URLParam(r, "puuid")

		err := db.ProfileUpdate(ctx, puuid)
		if err != nil {
			html.ServerError(w, r, nil, fmt.Errorf("updating summary: %w", err))
			return
		}

		w.Header().Set("HX-Refresh", "true")
		w.WriteHeader(http.StatusOK)
	}
}
