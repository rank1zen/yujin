package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

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

type matchV5ParticipantQuery struct {
	db pgxDB
}

func NewMatchV5ParticipantQuery(db pgxDB) IRecordQuery[MatchParticipantRecord] {
	return &matchV5ParticipantQuery{db: db}
}

func (q *matchV5ParticipantQuery) GetRecords(ctx context.Context, filters ...RecordFilter) ([]*MatchParticipantRecord, error) {
	return nil, nil
}

func (q *matchV5ParticipantQuery) CountRecords(ctx context.Context, filters ...RecordFilter) (int64, error) {
	return 0, nil
}

// GetMatchlist returns match IDs associated with a puuid
func getMatchlist(db pgxDB) func(context.Context, string) ([]string, error) {
	return func(ctx context.Context, puuid string) ([]string, error) {
		rows, _ := db.Query(ctx, `
                        SELECT match_id
                        FROM MatchParticipantRecords
                        WHERE puuid = $1
                `, puuid)

		defer rows.Close()
		records, err := pgx.CollectRows(rows, pgx.RowToStructByPos[string])
		if err != nil {
			return nil, fmt.Errorf("get matchlist: %w", err)
		}

		return records, nil
	}
}
