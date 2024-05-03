package database

import (
	"context"
	"fmt"
)

type MatchBanRecord struct {
	RecordId   string `db:"record_id"`
	MatchId    string `db:"match_id"`
	TeamId     int32  `db:"team_id"`
	ChampionId int    `db:"champion_id"`
	Turn       int    `db:"turn"`
}

type matchV5BanQuery struct {
	db pgxDB
}

func (q *matchV5BanQuery) GetRecords(ctx context.Context, filters ...RecordFilter) ([]*MatchBanRecord, error) {
	return nil, fmt.Errorf("not implemented")
}

func (q *matchV5BanQuery) CountRecords(ctx context.Context, flters ...RecordFilter) (int64, error) {
	return 0, fmt.Errorf("not implemented")
}

func (q *matchV5BanQuery) InsertRecords(ctx context.Context, records []MatchBanRecord) (int64, error) {
	return insertBulk[MatchBanRecord](ctx, q.db, "matchbanrecords", records)
}

func (q *matchV5BanQuery) DeleteRecords(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}
