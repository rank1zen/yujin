package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// FIXME: LOLOL

type LeagueRecord struct {
	RecordId   string    `db:"record_id"`
	RecordDate time.Time `db:"record_date"`
	SummonerId string    `db:"summoner_id"`
	LeagueId   string    `db:"league_id"`
	Tier       string    `db:"tier"`
	Rank       string    `db:"rank"`
	Lp         int32     `db:"league_points"`
	Wins       int32     `db:"number_wins"`
	Losses     int32     `db:"number_losses"`
}

type LeagueQuery interface {
	FetchAndInsertBySummoner(ctx context.Context, riot RiotClient, summonerId string) error
	GetRecent(ctx context.Context, summonerId string) (LeagueRecord, error)
}

type leagueQuery struct {
	db pgxDB
}

func NewLeagueQuery(db pgxDB) LeagueQuery {
	return &leagueQuery{db: db}
}

func (q *leagueQuery) FetchAndInsertBySummoner(ctx context.Context, riot RiotClient, summonerId string) error {
	league, err := riot.GetLeagueBySummoner(summonerId)
	if err != nil {
		return fmt.Errorf("fetch: %w", err)
	}

	_, err = q.db.Exec(ctx, `
	INSERT INTO LeagueRecords
	(summoner_id, league_id, tier, division, league_points, number_wins, number_losses)
	VALUES ($1, $2, $3, $4, $5, $6, $7);
	`, league.SummonerID, league, league.Tier, league, league.Wins, league.Losses)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	return nil
}

func (q *leagueQuery) GetRecent(ctx context.Context, summonerId string) (LeagueRecord, error) {
	rows, _ := q.db.Query(ctx, `
	FIXME PLEASE
	SELECT t1.* 
	FROM LeagueRecords AS t1
	JOIN (
		SELECT MAX(record_date) AS recent, puuid
		FROM LeagueRecords
		WHERE summoner_id = $1
	)
	`, summonerId)
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[LeagueRecord])
}
