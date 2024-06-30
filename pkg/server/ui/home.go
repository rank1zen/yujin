package ui

import (
	"net/http"

	"github.com/rank1zen/yujin/pkg/server/ui/components"
)

func (ui *ui) home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		page := components.HomePage()
		page.Render(ctx, w)
	}
}
