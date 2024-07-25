package ui

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rank1zen/yujin/pkg/database"
	"github.com/rank1zen/yujin/pkg/http/request"
	"github.com/rank1zen/yujin/pkg/server/ui/partials"
)

func profileMatchList(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		name, err := database.ParseRiotName(chi.URLParam(r, "name"))
		if err != nil {
			partials.ProfileMatchlistError().Render(ctx, w)
			return
		}

		page := request.QueryIntParam(r, "page", 10)

		matches, err := db.GetProfileMatchList(ctx, name, page)
		if err != nil {
			partials.ProfileMatchlistError().Render(ctx, w)
			return
		}

		partials.ProfileMatchlist(matches).Render(ctx, w)
	}
}
