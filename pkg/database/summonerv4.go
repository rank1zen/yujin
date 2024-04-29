package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/pkg/logging"
	"go.uber.org/zap"
)

type summonerV4Query struct {
	db pgDB
}

func NewSummonerV4Query(db pgDB) SummonerV4Query {
	return &summonerV4Query{db: db}
}

func (q *summonerV4Query) GetSummonerRecords(ctx context.Context, filters ...RecordFilter) ([]*SummonerRecord, error) {
	log := logging.FromContext(ctx)
	return getSummonerRecords(q.db, log)(ctx, filters...)
}

func (q *summonerV4Query) CountSummonerRecords(ctx context.Context, filters ...RecordFilter) (int64, error) {
	log := logging.FromContext(ctx)
	return countSummonerRecords(q.db, log)(ctx, filters...)
}

func (q *summonerV4Query) InsertSummonerRecords(ctx context.Context, records []*SummonerRecord) (int64, error) {
	return insertSummonerRecords(q.db)(ctx, records)
}

func (q *summonerV4Query) DeleteSummonerRecords(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}


func getSummonerRecords(db pgDB, log *zap.SugaredLogger) func(context.Context, ...RecordFilter) ([]*SummonerRecord, error) {
	return func(ctx context.Context, filters ...RecordFilter) ([]*SummonerRecord, error) {
		query := `
                        SELECT
                                record_id, record_date, puuid, account_id, id,
                                name, profile_icon_id, summoner_level, revision_date
                        FROM
                                summoner_records
                        WHERE 1=1
                `
		query, args := build(query, 0, filters...)

		log.Debugf("Read Records Query: %s", query)
		rows, _ := db.Query(ctx, query, args...)
		defer rows.Close()

		records, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[SummonerRecord])
		if err != nil {
			return nil, fmt.Errorf("select summmoner: %w", err)
		}

		return records, nil
	}
}

func countSummonerRecords(db pgDB, log *zap.SugaredLogger) func(context.Context, ...RecordFilter) (int64, error) {
	return func(ctx context.Context, filters ...RecordFilter) (int64, error) {
		query := `
                        SELECT
                                COUNT(*)
                        FROM
                                summoner_records
                        WHERE 1=1
                `
		query, args := build(query, 0, filters...)

		log.Debugf("Read Records Query: %s", query)

		var count int64
		err := db.QueryRow(ctx, query, args...).Scan(&count)
		if err != nil {
			return 0, fmt.Errorf("summoner count: %w", err)
		}

		return count, nil
	}
}

func insertSummonerRecords(db pgDB) func(context.Context, []*SummonerRecord) (int64, error) {
	return func(ctx context.Context, records []*SummonerRecord) (int64, error) {
		var rows [][]any
		for _, r := range records {
			rows = append(rows, []any{
				r.RecordDate, r.AccountId, r.ProfileIconId, r.RevisionDate, r.Name,
				r.SummonerId, r.Puuid, r.SummonerLevel,
			})
		}

		count, err := db.CopyFrom(ctx,
			pgx.Identifier{"summonerrecords"},
			[]string{
				"record_date", "account_id", "profile_icon_id", "revision_date", "name",
				"summoner_id", "puuid", "summoner_level",
			},
			pgx.CopyFromRows(rows))
		if err != nil {
			return 0, fmt.Errorf("insert summoner: %w", err)
		}

		return count, nil

	}
}
