package ui

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rank1zen/yujin/pkg/database"
	"github.com/rank1zen/yujin/pkg/server/ui/pages"
)

func aboutPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pages.HomePage("About", "V1.1", "aa", "aa").Render(ctx, w)
}

func profilePage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	name, err := database.ParseRiotName(chi.URLParam(r, "name"))
	if err != nil {
		pages.ProfileNotFoundErrorPage().Render(ctx, w)
		return
	}

	pages.ProfilePage(name).Render(ctx, w)
}
