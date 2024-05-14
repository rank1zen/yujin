package views

import (
	"fmt"
	"strconv"

	"github.com/rank1zen/yujin/pkg/components"
	"github.com/rank1zen/yujin/pkg/database"
)

type profilePage struct {
	summoner database.SummonerRecord

	// please make sure that these are all equal len
	matches      []database.MatchInfoRecord
	participants []database.MatchParticipantRecord
	//team
}

func (p profilePage) Matchlist() []components.MatchCardProps {
	matches := make([]components.MatchCardProps, len(p.matches))

	for i := range len(p.matches) {
		matches = append(matches, matchCard{
			match:       p.matches[i],
			participant: p.participants[i],
		})
	}

	return matches
}

func (p profilePage) Summoner() components.SummonerCardProps {
	return summonerCard{}
}

type summonerCard struct {
	summoner database.SummonerRecord
}

func (s summonerCard) Level() string { return fmt.Sprint(s.summoner.SummonerLevel) }
func (s summonerCard) Name() string  { return s.summoner.Name }
func (s summonerCard) ProfileIcon() string {
	// how do we do images?
	return ""
}

type matchCard struct {
	match       database.MatchInfoRecord
	participant database.MatchParticipantRecord
}

func (m matchCard) Win() bool            { return true }
func (m matchCard) GameDuration() string { return m.match.Duration.String() }
func (m matchCard) GamePatch() string    { return m.match.Patch }
func (m matchCard) GameDate() string     { return m.match.StartTs.String() }
func (m matchCard) KDA() string {
	return fmt.Sprintf("%d/%d/%d", m.participant.Kills, m.participant.Deaths, m.participant.Assists)
}

func (m matchCard) CS() string { return strconv.Itoa(m.participant.CreepScore) }
