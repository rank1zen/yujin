package ui

import (
	"net/http"

	"github.com/rank1zen/yujin/pkg/server/ui/pages"
)

func ErrorNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		w.WriteHeader(http.StatusNotFound)

		pages.ErrorNotFound().Render(ctx, w)
	}
}
