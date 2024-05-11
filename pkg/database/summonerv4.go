package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/pkg/logging"
)

type SummonerRecord struct {
	RecordId      string    `db:"record_id"`
	RecordDate    time.Time `db:"record_date"`
	Puuid         string    `db:"puuid"`
	AccountId     string    `db:"account_id"`
	SummonerId    string    `db:"id"`
	Name          string    `db:"name"`
	ProfileIconId int32     `db:"profile_icon_id"`
	SummonerLevel int32     `db:"summoner_level"`
	RevisionDate  int64     `db:"revision_date"`
}

type summonerV4Query struct {
	db pgxDB
}

func (q *summonerV4Query) GetRecords(ctx context.Context, filters ...RecordFilter) ([]*SummonerRecord, error) {
	log := logging.FromContext(ctx).Sugar()

	query := `
                SELECT
                        record_id, record_date, puuid, account_id, id,
                        name, profile_icon_id, summoner_level, revision_date
                FROM
                        summoner_records
                WHERE 1=1
        `
	// query, args := build(query, 0, filters...)
        args := []any{}

	log.Debugf("Read Records Query: %s", query)
	rows, _ := q.db.Query(ctx, query, args...)
	defer rows.Close()

	records, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[SummonerRecord])
	if err != nil {
		return nil, fmt.Errorf("get summmoner: %w", err)
	}

	return records, nil
}

func (q *summonerV4Query) CountRecords(ctx context.Context, filters ...RecordFilter) (int64, error) {
	log := logging.FromContext(ctx).Sugar()

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
	err := q.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("summoner count: %w", err)
	}

	return count, nil
}

func (q *summonerV4Query) InsertRecords(ctx context.Context, records []SummonerRecord) (int64, error) {
	var rows [][]any
	for _, r := range records {
		rows = append(rows, []any{
			r.RecordDate, r.AccountId, r.ProfileIconId, r.RevisionDate, r.Name,
			r.SummonerId, r.Puuid, r.SummonerLevel,
		})
	}

	count, err := q.db.CopyFrom(ctx,
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
