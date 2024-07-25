package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rank1zen/yujin/pkg/database"
	"github.com/rank1zen/yujin/pkg/server/ui/components"
)

func profileMatchlist(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_ = r.Header.Get("HX-Request")

		puuid := chi.URLParam(r, "puuid")

		pageParam := r.URL.Query().Get("page")

		var page int
		switch pageParam {
		case "":
			page = 0
		default:
			var err error
			page, err = strconv.Atoi(pageParam)
			if err != nil {
				page = 0
			}
		}



		matches, err := db.GetMatchHistory(ctx, puuid, page)
		if err != nil {
			components.ServerError(err).Render(ctx, w)
			return
		}

		components.MatchCardList(matches).Render(ctx, w)
	}
}
