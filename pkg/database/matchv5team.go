package database

import (
	"context"
	"fmt"
)

type MatchTeamRecord struct {
	RecordId  string `db:"record_id"`
	MatchId   string `db:"match_id"`
	TeamId    int32  `db:"team_id"`
	Win       bool   `db:"win"`
	Surrender bool   `db:"surrender"`
}

type matchV5TeamQuery struct {
	db pgxDB
}

func (q *matchV5TeamQuery) GetRecords(ctx context.Context, filters ...RecordFilter) ([]*MatchTeamRecord, error) {
	return nil, fmt.Errorf("not implemented")
}

func (q *matchV5TeamQuery) CountRecords(ctx context.Context, filters ...RecordFilter) (int64, error) {
	return 0, fmt.Errorf("not implemented")
}

func (q *matchV5TeamQuery) InsertRecords(ctx context.Context, records []MatchTeamRecord) (int64, error) {
	return insertBulk[MatchTeamRecord](ctx, q.db, "matchteamrecords", records)
}

func (q *matchV5TeamQuery) DeleteRecords(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}
