package database

import (
	"fmt"
	"time"
)

func (m ProfileMatch) GetGamePatch() string {
	return m.GamePatch[:5]
}

func (m ProfileMatch) GetGameDate() string {
	duration := time.Since(m.GameDate)

	switch {
	case duration < time.Minute:
		return fmt.Sprintf("%d seconds ago", int(duration.Seconds()))
	case duration < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	case duration < time.Hour*24:
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	case duration < time.Hour*24*7:
		return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
	default:
		return m.GameDate.Format("2006-01-02")
	}
}

func (m ProfileMatch) GetGameDuration() string {
	return fmt.Sprintf("%.0fm", m.GameDuration.Minutes())
}

type ProfileMatchList []*ProfileMatch

type ProfileMatchSummary struct {
	GameDate         time.Time     `db:"game_date"`
	GameDuration     time.Duration `db:"game_duration"`
	MatchId          string        `db:"match_id"`
	GamePatch        string        `db:"game_patch"`
	PlayerWin        bool          `db:"player_win"`
	PostGames        MatchSummonerPostGameList
	MatchItemEvents  []*MatchItemEvent
	MatchSkillEvents []*MatchSkillEvent
	MatchKillEvents  []*MatchKillEvent
}

type ProfileMatchSummaryList []*ProfileMatchSummary
