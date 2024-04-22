package summoner

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/internal/logging"
	"github.com/rank1zen/yujin/pkg/database"
)

type SummonerV4Query struct {
        db *database.DB
}

func NewSummonerV4Query(db *database.DB) *SummonerV4Query {
	return &SummonerV4Query{
		db: db,
	}
}

func (q *SummonerV4Query) CountSummonerRecords(ctx context.Context, f *SummonerRecordCountFilter) (int64, error) {
	var count int64
	err := q.db.Pgx.QueryRow(ctx, `
                        SELECT count(*) FROM summoner_records
                        WHERE
                                $1 = $2 AND record_date >= $3 AND record_date <= $4
                `,
                f.Field,
                f.Value,
                f.DateMin,
                f.DateMax,
        ).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("query error: %w", err)
	}

	return count, nil
}

func (q *SummonerV4Query) InsertBatchSummonerRecord(ctx context.Context, records []*SummonerRecordArg) ([]*string, error) {
        ids := make([]*string, len(records))
        var cur int;
        batch := &pgx.Batch{}

        scan := func (row pgx.Row) error {
                err := row.Scan(&ids[cur])
                if err != nil {
                        return err
                }

                cur++
                return nil
        }

        for _, r := range records {
                query := batch.Queue(`
                        INSERT INTO summoner_records
                        (record_date, account_id, profile_icon_id, revision_date, name, id, puuid, summoner_level)
                        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
                        RETURNING record_id
                `,
                r.RecordDate, r.AccountId, r.ProfileIconId, r.RevisionDate, r.Name, r.SummonerId,
                r.Puuid, r.SummonerLevel)

                query.QueryRow(scan)
        }
        err := q.db.Pgx.SendBatch(ctx, batch).Close()
        if err != nil {
                return nil, err
        }

        return ids, nil
}

func (q *SummonerV4Query) InsertSummonerRecord(ctx context.Context, a *SummonerRecordArg) (string, error) {
        log := logging.FromContext(ctx)
        log.Infof("Creating Summoner Record: %s, ", a.Name)

	var uuid string
	err := q.db.Pgx.QueryRow(ctx,
                `INSERT INTO summoner_records
                (record_date, account_id, profile_icon_id, revision_date, name, id, puuid, summoner_level)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
                RETURNING record_id`,
		a.RecordDate,
		a.AccountId,
		a.ProfileIconId,
		a.RevisionDate,
		a.Name,
		a.SummonerId,
		a.Puuid,
		a.SummonerLevel,
	).Scan(&uuid)
	if err != nil {
		return "", fmt.Errorf("query error: %w", err)
	}

        log.Infof("Created Summoner Record: %s", uuid)
	return uuid, nil
}

func (q *SummonerV4Query) ReadSummonerRecords(ctx context.Context, f *SummonerRecordFilter) ([]*SummonerRecordDatum, error) {
	rows, _ := q.db.Pgx.Query(ctx, `
                        SELECT
                                record_id, record_date, name, profile_icon_id, summoner_level, revision_date
                        FROM
                                summoner_records
                        WHERE
                                $1 = $2 AND record_date >= $3 AND record_date <= $4
                        ORDER BY record_date $5
                        LIMIT $6 OFFSET $7
                `,
                f.Field,
                f.Value,
                f.DateMin,
                f.DateMax,
                f.SortOrder,
                f.Limit,
                f.Offset,
        )
	defer rows.Close()

	records, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[SummonerRecordDatum])
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	return records, nil
}

func (q *SummonerV4Query) DeleteSummonerRecord(ctx context.Context, id string) error {
        log := logging.FromContext(ctx)
	_, err := q.db.Pgx.Exec(ctx,
                `DELETE FROM summoner_records
                WHERE record_id = $1`,
                id)
	if err != nil {
		return fmt.Errorf("query error: %w", err)
	}

        log.Infof("Deleted Record id: %s", id)
	return nil
}
