package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// SelectMatchBanRecords returns
func selectMatchBanRecords(db pgDB) func(context.Context) {
	return func(ctx context.Context) {
		qw := `
                SELECT
                        record_id, match_id, team_id,
                        champion_id, turn
                FROM
                        MatchBanRecords
                WHERE 1=1

                SELECT *
                FROM MatchTeamRecords
                INNER JOIN MatchRecords
                ON MatchRecords.id = 1
        `
		db.Query(ctx, qw)
	}
}

func insertMatchTeamRecords(db pgDB) func(context.Context, []*MatchTeamRecord) (int64, error) {
	return func(ctx context.Context, records []*MatchTeamRecord) (int64, error) {
		var rows [][]any
		for _, r := range records {
			rows = append(rows, []any{r.MatchId, r.TeamId, r.Win, r.Surrender})
		}

		count, err := db.CopyFrom(ctx,
			pgx.Identifier{"matchteamrecords"},
			[]string{"match_id", "team_id", "win", "surrender"},
			pgx.CopyFromRows(rows))
		if err != nil {
			return 0, fmt.Errorf("insert match team: %w", err)
		}

		return count, nil
	}
}

func insertMatchBanRecords(db pgDB) func(context.Context, []*MatchBanRecord) (int64, error) {
	return func(ctx context.Context, records []*MatchBanRecord) (int64, error) {
		var rows [][]any
		for _, r := range records {
			rows = append(rows, []any{r.MatchId, r.TeamId, r.ChampionId, r.Turn})
		}

		count, err := db.CopyFrom(ctx,
			pgx.Identifier{"matchbanrecords"},
			[]string{"match_id", "team_id", "champion_id", "turn"},
			pgx.CopyFromRows(rows))
		if err != nil {
			return 0, fmt.Errorf("insert match ban: %w", err)
		}

		return count, nil
	}
}

func insertMatchObjectiveRecords(db pgDB) func(context.Context, []*MatchObjectiveRecord) (int64, error) {
	return func(ctx context.Context, records []*MatchObjectiveRecord) (int64, error) {
		var rows [][]any
		for _, r := range records {
			rows = append(rows, []any{r.MatchId, r.TeamId, r.First, r.Kills, r.Name})
		}

		count, err := db.CopyFrom(ctx,
			pgx.Identifier{"matchobjectiverecords"},
			[]string{"match_id", "team_id", "first", "kills", "name"},
			pgx.CopyFromRows(rows))
		if err != nil {
			return 0, fmt.Errorf("insert match obj: %w", err)
		}

		return count, nil

	}
}
