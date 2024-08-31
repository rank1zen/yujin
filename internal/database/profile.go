package database

import (
	"time"
)

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
