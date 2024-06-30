package ui

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rank1zen/yujin/pkg/server/ui/components"
)

func (ui *ui) profile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		puuid := chi.URLParam(r, "puuid")
		matches, err := ui.db.ProfilePage(ctx, puuid)

		if err != nil {
			components.ServerError(err).Render(ctx, w)
			return
		}

		components.ProfilePage(*matches).Render(ctx, w)
	}
}
