package database

import (
	"time"
)

// A MatchRecord represents a row of MatchRecords
type MatchRecord struct {
	RecordId   string        `db:"record_id"`
	RecordDate time.Time     `db:"record_date"`
	MatchId    string        `db:"match_id"`
	StartTs    time.Time     `db:"start_ts"`
	Duration   time.Duration `db:"duration"`
	Surrender  bool          `db:"surrender"`
	Patch      string        `db:"patch"`
}

// A FullMatchRecord represents a collection of rows associated to one particular match
type FullMatchRecord struct {
	Metadata     *MatchRecord
	Teams        []*MatchTeamRecord
	Objectives   []*MatchObjectiveRecord
	Bans         []*MatchBanRecord
	Participants []*MatchParticipantRecord
}

type MatchTeamRecord struct {
	RecordId  string `db:"record_id"`
	MatchId   string `db:"match_id"`
	TeamId    int32  `db:"team_id"`
	Win       bool   `db:"win"`
	Surrender bool   `db:"surrender"`
}

type MatchBanRecord struct {
	RecordId   string `db:"record_id"`
	MatchId    string `db:"match_id"`
	TeamId     int32  `db:"team_id"`
	ChampionId int    `db:"champion_id"`
	Turn       int    `db:"turn"`
}

type MatchObjectiveRecord struct {
	RecordId string `db:"record_id"`
	MatchId  string `db:"match_id"`
	TeamId   int32  `db:"team_id"`
	Name     string `db:"name"`
	First    bool   `db:"first"`
	Kills    int    `db:"kills"`
}

// MatchParticipantRecord represents a record of a participant in a match
// stored in the database
type MatchParticipantRecord struct {
	RecordId string `db:"record_id"`
	MatchId  string `db:"match_id"`
	Puuid    string `db:"puuid"`

	ParticipantId int    `db:"participant_id"`
	TeamId        int    `db:"team_id"`
	SummonerName  string `db:"summoner_name"`
	SummonerLevel int    `db:"summoner_level"`
	Position      string `db:"position"`
	ChampId       int    `db:"champion_id"`
	ChampName     string `db:"champion_name"`
	ChampLevel    int    `db:"champion_level"`

	Kills      int `db:"kills"`
	Deaths     int `db:"deaths"`
	Assists    int `db:"assists"`
	CreepScore int `db:"creep_score"`
	GoldEarned int `db:"gold_earned"`

	VisionScore        int `db:"VisionScore"`
	WardsPlaced        int `db:"WardsPlaced"`
	ControlWardsPlaced int `db:"ControlWardsPlaced"`

	FirstBloodAssist bool `db:"FirstBloodAssist"`
	FirstTowerAssist bool `db:"FirstTowerAssist"`
	TurretTakeDowns  int  `db:"TurretTakeDowns"`

	PhysicalDamageDealtToChampions int `db:"PhysicalDamageDealtToChampions"`
	MagicDamageDealtToChampions    int `db:"MagicDamageDealtToChampions"`
	TrueDamageDealtToChampions     int `db:"TrueDamageDealtToChampions"`
	TotalDamageDealtToChampions    int `db:"TotalDamageDealtToChampions"`
	TotalDamageTaken               int `db:"TotalDamageTaken"`
	TotalHealsOnTeammates          int `db:"TotalHealsOnTeammates"`
}
