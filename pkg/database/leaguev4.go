package database

import (
	"context"
	"fmt"
	"time"
)

type LeagueRecord struct {
	RecordId     string    `db:"record_id"`
	RecordDate   time.Time `db:"record_date"`
	LeagueId     string    `db:"league_id"`
	QueueType    string    `db:"queue_type"`
	SummonerId   string    `db:"summoner_id"`
	Tier         string    `db:"tier"`
	Rank         string    `db:"rank"`
	LeaguePoints int32     `db:"league_points"`
	Wins         int32     `db:"wins"`
	Losses       int32     `db:"losses"`
}

type leagueV4Query struct {
	db pgxDB
}

func (q *leagueV4Query) GetRecords(ctx context.Context, filters ...RecordFilter) ([]*LeagueRecord, error) {
	return nil, fmt.Errorf("not implemented")
}

func (q *leagueV4Query) CountRecords(ctx context.Context, flters ...RecordFilter) (int64, error) {
	return 0, fmt.Errorf("not implemented")
}

func (q *leagueV4Query) InsertRecords(ctx context.Context, records []LeagueRecord) (int64, error) {
	return 0, fmt.Errorf("not implemented")
}

func (q *leagueV4Query) DeleteRecords(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}
