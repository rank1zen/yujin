package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type MatchRecordArg struct {
	RecordDate time.Time
	MatchId    string
	StartTs    int
	Duration   int64
	Surrender  bool
	Patch      string
	Puuids     []string

	//RedSide  Team
	//BlueSide Team
	Players []MatchParticipantRecordArg
}

type MatchRecord struct {
	RecordId   string    `db:"record_id"`
	RecordDate time.Time `db:"record_date"`
	MatchId    string    `db:"match_id"`
	StartTs    int       `db:"start_ts"`
	Duration   int64     `db:"duration"`
	Surrender  bool      `db:"surrender"`
	Patch      string    `db:"patch"`
}

type MatchTeamRecordArg struct {
	RecordDate time.Time
	MatchId    string
	TeamId     int
	Bans       TeamBan
	Objective  TeamObjective
}

type MatchTeamRecord struct {
	RecordId   string          `db:"record_id"`
	RecordDate time.Time       `db:"record_date"`
	MatchId    string          `db:"match_id"`
	TeamId     int32           `db:"team_id"`
	Bans       []TeamBan       `db:"bans"`
	Objectives []TeamObjective `db:"objectives"`
}

type MatchParticipantRecordArg struct {
	Puuid         string
	SummonerLevel int
	SummonerName  string
	Id            int

	Win        bool
	Kills      int
	Deaths     int
	Assists    int
	CS         int
	Gold       int
	Position   string
	ChampId    int
	ChampName  int
	ChampLevel int

	Runes Perks

	VisionScore        int
	WardsPlaced        int
	ControlWardsPlaced int

	FirstBloodAssist bool
	FirstTowerAssist bool
	TurretTakeDowns  int

	PhysicalDamageDealtToChampions int
	MagicDamageDealtToChampions    int
	TrueDamageDealtToChampions     int
	TotalDamageDealtToChampions    int
	TotalDamageTaken               int
	TotalHealsOnTeammates          int

	Item0 int
	Item1 int
	Item2 int
	Item3 int
	Item4 int
	Item5 int
	Item6 int
}

type MatchParticipantRecord struct {
	SummonerName string `db:"record_date"`

	Win        bool   `db:"record_date"`
	Kills      int    `db:"record_date"`
	Deaths     int    `db:"record_date"`
	Assists    int    `db:"record_date"`
	CS         int    `db:"record_date"`
	Gold       int    `db:"record_date"`
	Position   string `db:"record_date"`
	ChampId    int    `db:"record_date"`
	ChampName  int    `db:"record_date"`
	ChampLevel int    `db:"record_date"`

	Runes Perks `db:"record_date"`

	VisionScore        int `db:"record_date"`
	WardsPlaced        int `db:"record_date"`
	ControlWardsPlaced int `db:"record_date"`

	FirstBloodAssist bool `db:"record_date"`
	FirstTowerAssist bool `db:"record_date"`
	TurretTakeDowns  int  `db:"record_date"`

	PhysicalDamageDealtToChampions int `db:"record_date"`
	MagicDamageDealtToChampions    int `db:"record_date"`
	TrueDamageDealtToChampions     int `db:"record_date"`
	TotalDamageDealtToChampions    int `db:"record_date"`
	TotalDamageTaken               int `db:"record_date"`
	TotalHealsOnTeammates          int `db:"record_date"`

	Item0 int `db:"record_date"`
	Item1 int `db:"record_date"`
	Item2 int `db:"record_date"`
	Item3 int `db:"record_date"`
	Item4 int `db:"record_date"`
	Item5 int `db:"record_date"`
	Item6 int `db:"record_date"`
}

type TeamBan struct {
	ChampionId int
	PickTurn   int
}

type TeamObjective struct {
	Name  string
	First bool
	Kills int
}

func RegisterTeamBanType(ctx context.Context, conn *pgx.Conn) error {
	dt, err := conn.LoadType(ctx, "team_champion_ban")
	if err != nil {
		return err
	}

	conn.TypeMap().RegisterType(dt)

	return nil
}

func RegisterTeamObjectiveType(ctx context.Context, conn *pgx.Conn) error {
	dt, err := conn.LoadType(ctx, "team_objective")
	if err != nil {
		return err
	}

	conn.TypeMap().RegisterType(dt)

	return nil
}

func (t TeamBan) IsNull() bool {
	return false
}

type Perks struct {
	Styles    []PerkStyle
	StatPerks PerkStats
}

type PerkStyle struct {
	Desc       string
	Style      int
	Selections []PerkStyleSelection
}

type PerkStyleSelection struct {
	Perk int
	Var1 int
	Var2 int
	Var3 int
}

type PerkStats struct {
	Defense int32
	Flex    int32
	Offense int32
}
