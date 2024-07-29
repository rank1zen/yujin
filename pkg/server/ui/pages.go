package ui

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rank1zen/yujin/pkg/database"
	"github.com/rank1zen/yujin/pkg/server/ui/pages"
)

func aboutPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pages.AboutPage("About", "V1.1", "aa", "aa").Render(ctx, w)
}

func profilePage(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		nameParam := chi.URLParam(r, "name")

		name, err := database.ParseRiotName(nameParam)
		if err != nil {
			pages.ProfileNotFoundErrorPage(name).Render(ctx, w)
			return
		}

		sum, err := db.GetProfileSummary(ctx, name)
		if err != nil {
			pages.ProfileNotFoundErrorPage(name).Render(ctx, w)
			log.Printf("hi: %v", err)
			return
		}

		pages.ProfilePage(nameParam, sum).Render(ctx, w)
	}
}
