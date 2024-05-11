package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type MatchRecord struct {
	RecordId   string        `db:"record_id"`
	RecordDate time.Time     `db:"record_date"`
	MatchId    string        `db:"match_id"`
	StartTs    time.Time     `db:"start_ts"`
	Duration   time.Duration `db:"duration"`
	Surrender  bool          `db:"surrender"`
	Patch      string        `db:"patch"`
}

type matchV5Query struct {
	db pgxDB
}

func (q *matchV5Query) GetRecords(ctx context.Context, filters ...RecordFilter) ([]*MatchRecord, error) {
	rows, _ := q.db.Query(ctx, `
                SELECT
                        record_id, record_date, match_id, start_ts, duration, surrender, patch
                FROM
                        MatchRecords`)

	defer rows.Close()
	records, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[MatchRecord])
	if err != nil {
		return nil, fmt.Errorf("select match: %w", err)
	}

	return records, nil
}

func (q *matchV5Query) CountRecords(ctx context.Context, filters ...RecordFilter) (int64, error) {
	var count int64
	err := q.db.QueryRow(ctx, `SELECT COUNT(*) FROM MatchRecords`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count match: %w", err)
	}

	return count, nil
}
