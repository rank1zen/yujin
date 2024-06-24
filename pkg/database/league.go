package database

import (
	"context"
	"time"
)

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

type leagueQuery struct {
	db pgxDB
}

// TODO: implement
func (c *service) FetchAndInsertBySummoner(ctx context.Context, summonerId string) error {
	return nil
}

// TODO: implement
func (c *service) GetRecentBySummoner(ctx context.Context, summonerId string) (*LeagueRecord, error) {
	return nil, nil
}
