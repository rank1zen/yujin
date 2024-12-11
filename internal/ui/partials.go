package ui

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/database"
	"github.com/rank1zen/yujin/internal/http/request"
	"github.com/rank1zen/yujin/internal/http/response/html"
	"github.com/rank1zen/yujin/internal/riot"
	"github.com/rank1zen/yujin/internal/ui/partials"
)

func checkHTMX(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !request.IsHTMX(r) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	})
}

func profileShowMatchList(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func profileShowChampStats(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func profileShowLiveGame(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}
