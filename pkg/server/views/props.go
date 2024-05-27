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

func (p profilePage) ProfileCard() components.ProfileCardProps {
	return profileCard{}
}

type profileCard struct {
	summoner database.SummonerRecord
	rank     database.LeagueRecord
}

func (p profileCard) IsRanked() bool      { return true }
func (p profileCard) LP() string          { return strconv.Itoa(int(p.rank.Lp)) }
func (p profileCard) Level() string       { return strconv.Itoa(int(p.summoner.SummonerLevel)) }
func (p profileCard) Losses() string      { return strconv.Itoa(int(p.rank.Losses)) }
func (p profileCard) Name() string        { return "FIXME: WHAT IS THE new RIOT NAMEs" }
func (p profileCard) ProfileIcon() string { return "FIXME: IMAGES?" }
func (p profileCard) Rank() string        { return p.rank.Rank }
func (p profileCard) Tier() string        { return p.rank.Tier }
func (p profileCard) Wins() string        { return strconv.Itoa(int(p.rank.Wins)) }

type matchCard struct {
	match       database.MatchInfoRecord
	participant database.MatchParticipantRecord
}

func (m matchCard) Win() bool            { return true }
func (m matchCard) GameDuration() string { return m.match.Duration.String() }
func (m matchCard) GamePatch() string    { return m.match.Patch }
func (m matchCard) GameDate() string     { return m.match.StartTs.String() }
func (m matchCard) CS() string           { return strconv.Itoa(m.participant.CreepScore) }
func (m matchCard) KDA() string {
	return fmt.Sprintf("%d/%d/%d", m.participant.Kills, m.participant.Deaths, m.participant.Assists)
}
