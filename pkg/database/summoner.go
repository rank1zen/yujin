package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type SummonerRecord struct {
	RecordId      string    `db:"record_id"`
	RecordDate    time.Time `db:"record_date"`
	Puuid         string    `db:"puuid"`
	AccountId     string    `db:"account_id"`
	SummonerId    string    `db:"summoner_id"`
	ProfileIconId int32     `db:"profile_icon_id"`
	SummonerLevel int32     `db:"summoner_level"`
	RevisionDate  time.Time `db:"revision_date"`
}

// FIXME: A simple sorting thing for summoners
// Which summoner (puuid), from dates A to B, Sort by Date ASC or DESC
// FIXME: pagination yes or no?
type SummonerRecordFilteR struct {
	Puuid   string
	DateMin time.Time
	DateMax time.Time
	DateAsc bool
}

type SummonerQuery interface {
	FetchAndInsert(ctx context.Context, gc RiotClient, puuid string) error
	GetRecent(ctx context.Context, puuid string) (SummonerRecord, error)

	// FIXME: these filters are rubbish mate
	GetRecords(ctx context.Context, filter SummonerRecordFilteR) ([]SummonerRecord, error)
	CountRecords(ctx context.Context, filter SummonerRecordFilteR) (int64, error)
}

type summonerQuery struct {
	db pgxDB
}

func NewSummonerQuery(db pgxDB) SummonerQuery {
	return &summonerQuery{db: db}
}

func (q *summonerQuery) FetchAndInsert(ctx context.Context, gc RiotClient, puuid string) error {
	summ, err := gc.GetSummoner(puuid)
	if err != nil {
		return fmt.Errorf("fetch: %w", err)
	}

	revDate := time.Unix(int64(summ.RevisionDate), 0)

	_, err = q.db.Exec(ctx, `
	INSERT INTO SummonerRecords
	(account_id, summoner_id, puuid, profile_icon_id, revision_date, summoner_level)
	VALUES ($1, $2, $3, $4, $5, $6)
        `, summ.AccountID, summ.ID, summ.PUUID, summ.ProfileIconID, revDate, summ.SummonerLevel)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	return nil
}

func (q *summonerQuery) GetRecent(ctx context.Context, puuid string) (SummonerRecord, error) {
	// HACK: Check that this query is actually good
	rows, _ := q.db.Query(ctx, `
	SELECT t1.*
	FROM SummonerRecords AS t1
	JOIN (
		SELECT MAX(record_date) AS recent, puuid
		FROM SummonerRecords
		WHERE puuid = $1
		GROUP BY puuid
	) AS t2 ON t2.puuid = t1.puuid AND t2.recent = t1.record_date
	WHERE t1.puuid = $1;
	`, puuid)

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[SummonerRecord])
}

func (q *summonerQuery) GetRecords(ctx context.Context, filter SummonerRecordFilteR) ([]SummonerRecord, error) {
	return nil, fmt.Errorf("not implemented")
}

func (q *summonerQuery) CountRecords(ctx context.Context, filter SummonerRecordFilteR) (int64, error) {
	return 0, fmt.Errorf("not implemented")
}
