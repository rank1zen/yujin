package database

import (
	"context"
	"fmt"
)

type MatchObjectiveRecord struct {
	RecordId string `db:"record_id"`
	MatchId  string `db:"match_id"`
	TeamId   int32  `db:"team_id"`
	Name     string `db:"name"`
	First    bool   `db:"first"`
	Kills    int    `db:"kills"`
}

type matchV5ObjQuery struct {
	db pgxDB
}

func (q *matchV5ObjQuery) GetRecords(ctx context.Context, filters ...RecordFilter) ([]*MatchObjectiveRecord, error) {
	return nil, fmt.Errorf("not implemented")
}

func (q *matchV5ObjQuery) CountRecords(ctx context.Context, filters ...RecordFilter) (int64, error) {
	return 0, fmt.Errorf("not implemented")
}
