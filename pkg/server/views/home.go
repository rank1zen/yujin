package views

import (
	"net/http"

	"github.com/rank1zen/yujin/pkg/components"
)

func (s *handler) home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		page := components.HomePage()
		page.Render(ctx, w)
	}
}
