package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/pkg/logging"
)

// FIXME: A bunch of things need to get fixed here

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

type SummonerQuery interface {
	FetchAndInsert(ctx context.Context, gc RiotClient, puuid string) error
	GetRecent(ctx context.Context, puuid string) (SummonerRecord, error)
	GetRecords(ctx context.Context, filters ...RecordFilter) ([]SummonerRecord, error)
	CountRecords(ctx context.Context, filters ...RecordFilter) (int64, error)
}

type summonerQuery struct {
	db pgxDB
}

func NewSummonerQuery(db pgxDB) SummonerQuery {
	return &summonerQuery{db: db}
}

func (q *summonerQuery) FetchAndInsert(ctx context.Context, gc RiotClient, puuid string) error {
	timestamp := time.Now()
	summ, err := gc.GetSummoner(puuid)
	if err != nil {
		return fmt.Errorf("GetByPUUID: %w", err)
	}

	_, err = q.db.Exec(ctx, `
                INSERT INTO SummonerRecords
                (record_date, account_id, id, name, puuid, profile_icon_id, revision_date, summoner_level)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        `,
		timestamp, summ.AccountID, summ.ID, summ.Name, summ.PUUID,
		summ.ProfileIconID, summ.RevisionDate, summ.SummonerLevel)
	if err != nil {
		return fmt.Errorf("insert summoner: %w", err)
	}

	return nil
}

func (q *summonerQuery) GetRecent(ctx context.Context, puuid string) (SummonerRecord, error) {
	// TODO: Complete this function
	return SummonerRecord{}, nil
}

func (q *summonerQuery) GetRecords(ctx context.Context, filters ...RecordFilter) ([]SummonerRecord, error) {
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

	records, err := pgx.CollectRows(rows, pgx.RowToStructByName[SummonerRecord])
	if err != nil {
		return nil, fmt.Errorf("get summmoner: %w", err)
	}

	return records, nil
}

func (q *summonerQuery) CountRecords(ctx context.Context, filters ...RecordFilter) (int64, error) {
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
