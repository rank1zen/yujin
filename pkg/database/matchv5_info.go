package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func BulkIns[T any](db pgDB, table string) func(context.Context, []T) (int64, error) {
	return func(ctx context.Context, a []T) (int64, error) {
		fields, rows, err := ExtractStructSlice(a)
		count, err := db.CopyFrom(ctx, pgx.Identifier{table}, fields, pgx.CopyFromRows(rows))
                if err != nil {
                        return 0, err
                }

                return count, nil
	}
}

// InsertMatchRecords inserts a list of match info records
func insertMatchRecords(db pgDB) func(ctx context.Context, records []MatchRecord) (int64, error) {
        f := BulkIns[MatchRecord](db, "matchrecords")
	return func(ctx context.Context, records []MatchRecord) (int64, error) {
                return f(ctx, records)
		var rows [][]any
		for _, r := range records {
			rows = append(rows, []any{r.RecordDate, r.MatchId, r.StartTs, r.Duration, r.Surrender, r.Patch})
		}

		count, err := db.CopyFrom(ctx,
			pgx.Identifier{"matchrecords"},
			[]string{"record_date", "match_id", "start_ts", "duration", "surrender", "patch"},
			pgx.CopyFromRows(rows))
		if err != nil {
			return 0, fmt.Errorf("insert match: %w", err)
		}

		return count, nil
	}
}

// SelectMatchRecords returns a list of match info records that satisfies some filter
func selectMatchRecords(db pgDB) func(context.Context, ...RecordFilter) ([]*MatchRecord, error) {
	return func(ctx context.Context, f ...RecordFilter) ([]*MatchRecord, error) {

		rows, _ := db.Query(ctx, "")

		defer rows.Close()
		records, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[MatchRecord])
		if err != nil {
			return nil, fmt.Errorf("select match: %w", err)
		}

		return records, nil
	}
}

// CountMatchRecords returns the count of match info records that satisfies some filter
func countMatchRecords(db pgDB) func(context.Context) (int64, error) {
	return func(ctx context.Context) (int64, error) {
		var count int64
		err := db.QueryRow(ctx, `
                        SELECT COUNT(*)
                        FROM MatchRecords
                `).Scan(&count)
		if err != nil {
			return 0, fmt.Errorf("count match: %w", err)
		}

		return count, nil
	}
}
