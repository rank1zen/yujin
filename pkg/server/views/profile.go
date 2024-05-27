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
		ctx := r.Context()

		r.ParseForm()

		puuid := r.FormValue("puuid")

		// we should be getting the most recent one
		summoner, err := s.db.Summoner().GetRecent(ctx, puuid)
		if err != nil {
			// TODO: how to handle errors?
		}

		soloq, err := s.db.League().GetRecentBySummoner(ctx, summoner.SummonerId)
		if err != nil {
			// TODO: how to handle errors?
		}

		_, err = s.db.Match().GetMatchlist(ctx, puuid)
		if err != nil {
			// TODO: how to handle errors?
		}

		comp := components.ProfilePage(components.ProfilePageProps{
			Profile: profileCard{
				summoner: summoner,
				rank: soloq,
			},
			Matchlist: make([]components.MatchCardProps, 0),
		})

		comp.Render(ctx, w)
	}
}
