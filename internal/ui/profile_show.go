package ui

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rank1zen/yujin/internal/database"
	"github.com/rank1zen/yujin/internal/http/response/html"
	"github.com/rank1zen/yujin/internal/ui/pages"
)

func profileShow(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}
