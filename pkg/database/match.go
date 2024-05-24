package database

import (
	"context"
	"fmt"
	"time"

	"github.com/KnutZuidema/golio/riot/lol"
	"github.com/jackc/pgx/v5"
)

// FIXME: LOLOL

type MatchRecord struct {
	info MatchInfoRecord
}

type MatchInfoRecord struct {
	RecordId   string        `db:"record_id"`
	RecordDate time.Time     `db:"record_date"`
	MatchId    string        `db:"match_id"`
	Patch      string        `db:"patch"`
	Duration   time.Duration `db:"duration"`
}

type MatchObjectiveRecord struct {
	RecordId string `db:"record_id"`
	MatchId  string `db:"match_id"`
	TeamId   int32  `db:"team_id"`
	Name     string `db:"name"`
	First    bool   `db:"first"`
	Kills    int    `db:"kills"`
}

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

type MatchBanRecord struct {
	RecordId   string `db:"record_id"`
	MatchId    string `db:"match_id"`
	TeamId     int32  `db:"team_id"`
	ChampionId int    `db:"champion_id"`
	Turn       int    `db:"turn"`
}

type MatchTeamRecord struct {
	RecordId  string `db:"record_id"`
	MatchId   string `db:"match_id"`
	TeamId    int32  `db:"team_id"`
	Win       bool   `db:"win"`
	Surrender bool   `db:"surrender"`
}

type MatchQuery interface {
	FetchAndInsert(ctx context.Context, riot RiotClient, puuid string) error

	// Fetch all the matches on the account  
	FetchAndInsertAll(ctx context.Context, riot RiotClient, puuid string) error

	// Get most recent matches 
	GetMatchlist(ctx context.Context, puuid string) ([]MatchRecord, error)

	// TODO: Implement these
	// GetBanRecords()
	// CountBanRecords()
	// GetObjectiveRecords()
	// CountObjectiveRecords()
}

type matchQuery struct {
	db pgxDB
}

func NewMatchQuery(db pgxDB) MatchQuery {
	return &matchQuery{db: db}
}

func (q *matchQuery) FetchAndInsert(ctx context.Context, riot RiotClient, puuid string) error {
	ids, err := riot.GetMatchlist(puuid)
	if err != nil {
		return fmt.Errorf("failed to fetch ids: %w", err)
	}

	rows, _ := q.db.Query(ctx, `
	SELECT match_id FROM MatchRecords WHERE match_id IN ANY($1)
	`, ids)

	newIDs, err := pgx.CollectRows(rows, pgx.RowToStructByName[string])
	if err != nil {
		return fmt.Errorf("failed to check db for new ids: %w", err)
	}

	err = pgx.BeginFunc(ctx, q.db, fetchAndInsertMatches(ctx, riot, newIDs))
	if err != nil {
		return fmt.Errorf("failed to fetch and insert: %w", err)
	}
	
	return nil
}

func (q *matchQuery) FetchAndInsertAll(ctx context.Context, riot RiotClient, puuid string) error {
	return fmt.Errorf("not implemented")
}

// fetchAndInsertMatches returns a function to execute in transaction.
// Queries Riot for matches and inserts.
func fetchAndInsertMatches(ctx context.Context, riot RiotClient, ids []string) func(pgx.Tx) error {
	return func(tx pgx.Tx) error {
		for _, id := range ids {
			match, err := riot.GetMatch(id)
			if err != nil {
				return err
			}

			err = insertFullMatch(ctx, tx, match)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

// TODO: Add everything
func insertFullMatch(ctx context.Context, db pgx.Tx, m *lol.Match) error {
	matchID := m.Metadata.MatchID

	_, err := db.Exec(ctx, `
	INSERT INTO MatchInfoRecords
	(record_date, match_id, duration, patch)
	VALUES ($1, $2, $3, $4)
	`, m.Info.GameStartTimestamp, matchID, m.Info.GameDuration, m.Info.GameVersion)
	if err != nil {
		return fmt.Errorf("MatchInfo: %w", err)
	}

	for _, p := range m.Info.Participants {
		_, err := db.Exec(ctx, `
		INSERT INTO MatchParticipantRecords
		(match_id, puuid)
		VALUES ($1, $2)
		`, matchID, p.PUUID)
		if err != nil {
			return fmt.Errorf("MatchParticipant: %w", err)
		}
	}

	for _, t := range m.Info.Teams {
		_, err := db.Exec(ctx, `
		INSERT INTO MatchTeamRecords
		(match_id, team_id)
		`, matchID, t.TeamID)
		if err != nil {
			return fmt.Errorf("MatchTeam: %w", err)
		}
	}

	return nil
}

func (q *matchQuery) GetMatchlist(ctx context.Context, puuid string) ([]MatchRecord, error) {
	rows, _ := q.db.Query(ctx, `
	SELECT i.record_date, i.match_id, i.duration, i.patch
	FROM MatchInfoRecords AS i
		JOIN MatchParticipantRecords AS p ON i.match_id = p.match_id
	WHERE p.puuid = $1
	`, puuid)

	_, err := pgx.CollectRows(rows, pgx.RowToStructByName[MatchInfoRecord])
	if err != nil {
		return nil, nil
	}


	return nil, nil
}
