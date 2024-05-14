package views

import (
	"net/http"

	"github.com/rank1zen/yujin/pkg/components"
	"github.com/rank1zen/yujin/pkg/database"
)

type profilesHandler struct {
	db database.DB
}

func (s *profilesHandler) profile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// FIXME: error checking

		ctx := r.Context()

		r.ParseForm()

		// we should be getting the most recent one
		summoner, err := s.db.Summoner().GetRecent(ctx, "")

		// TODO: we should be getting games for this summoner
		matches, err := s.db.Match().GetMatchlist(ctx, "")

		comp := components.ProfilePage(profilePage{
			summoner: summoner[0],
		})
		comp.Render(ctx, w)
	}
}
