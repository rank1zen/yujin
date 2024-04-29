package database

import (
	"context"
	"time"
)

type leagueV4Query struct {
	db *DB
}

func NewLeagueV4Query(db *DB) LeagueV4Query {
	return &leagueV4Query{db: db}
}

// LeagueRecords represents a record of a league entry
type LeagueRecords struct {
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

// GetLeagueRecords returns a list of LeagueRecords satisfying something
func (q *leagueV4Query) GetLeagueRecords(ctx context.Context) {

}

// InsertLeagueRecords inserts a LeagueRecord
func (q *leagueV4Query) InsertLeagueRecords(ctx context.Context) {

}

func (q *leagueV4Query) CountLeagueRecords(ctx context.Context) {

}
