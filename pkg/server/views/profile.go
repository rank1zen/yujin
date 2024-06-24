package views

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rank1zen/yujin/pkg/components"

	"github.com/rank1zen/yujin/pkg/database"
)

func (s *handler) profile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		r.ParseForm()

		puuid := r.FormValue("puuid")
		summoner, err := s.db.GetMatchHistory(ctx, puuid)
		if err != nil {
			return // TODO: how to handle errors?
		}

		soloq, err := s.db.League().GetRecentBySummoner(ctx, summoner.SummonerId)
		if err != nil {
			return // TODO: how to handle errors?
		}

		matches, err := s.db.Match().GetMatchlist(ctx, puuid)
		if err != nil {
			return // TODO: how to handle errors?
		}

		ad := make([]components.MatchCardProps, len(matches))
		for i, m := range matches {
			ad[i] = matchCard{
				match: m,
			}
		}

		props := components.ProfilePageProps{
			Profile: ProfileCard{
				Summoner:          summoner,
				SummonerSoloqRank: soloq,
			},
			Matchlist: ad,
		}
		components.ProfilePage(props).Render(ctx, w)
	})
}

type ProfileCard struct {
	Summoner          database.SummonerRecord
	SummonerSoloqRank database.LeagueRecord
}

func (p ProfileCard) IsRanked() bool { return true }
func (p ProfileCard) LP() string     { return strconv.Itoa(int(p.SummonerSoloqRank.Lp)) }
func (p ProfileCard) Level() string  { return strconv.Itoa(int(p.Summoner.SummonerLevel)) }
func (p ProfileCard) Losses() string { return strconv.Itoa(int(p.SummonerSoloqRank.Losses)) }
func (p ProfileCard) Name() string   { return "Summoner Name" } // FIXME
func (p ProfileCard) ProfileIcon() string {
	return "https://static.bigbrain.gg/assets/lol/riot_static/14.10.1/img/champion/Jhin.png"
}
func (p ProfileCard) Rank() string { return p.SummonerSoloqRank.Rank }
func (p ProfileCard) Tier() string { return p.SummonerSoloqRank.Tier }
func (p ProfileCard) Wins() string { return strconv.Itoa(int(p.SummonerSoloqRank.Wins)) }

type matchCard struct {
	match       database.MatchInfoRecord
	participant database.MatchParticipantRecord
}

func (m matchCard) Win() bool            { return true }
func (m matchCard) GameDuration() string { return m.match.Duration.String() }
func (m matchCard) GamePatch() string    { return m.match.Patch }
func (m matchCard) GameDate() string     { return m.match.RecordDate.String() }
func (m matchCard) CS() string           { return strconv.Itoa(m.participant.CreepScore) }
func (m matchCard) KDA() string {
	return fmt.Sprintf("%d / %d / %d", m.participant.Kills, m.participant.Deaths, m.participant.Assists)
}
