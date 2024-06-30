package database

import (
	"context"
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

type SummonerProfileResponse struct {
	Name          string
	Level         int
	ProfileIconID int
	Wins          int
	Losses        int
	Rank          *string
	Tier          *string
	LP            *int
}

func (db *DB) countSummonerRecords(ctx context.Context, puuid string) (int64, error) {
	var c int64
	err := db.pool.QueryRow(ctx, `SELECT count(*) FROM summoner_records WHERE puuid = $1`, puuid).Scan(&c)
	if err != nil {
		return 0, err
	}

	return c, nil
}

func (db *DB) newestSummonerRecord(ctx context.Context, puuid string) (*SummonerRecord, error) {
	rows, _ := db.pool.Query(ctx, `SELECT * FROM summoner_records_newest WHERE puuid = $1`, puuid)

	record, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[SummonerRecord])
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (db *DB) UpdateSummoner(ctx context.Context, puuid string) (*SummonerRecord, error) {
	m, err := db.riot.GetSummoner(ctx, puuid)
	if err != nil {
		return nil, err
	}

	ts := time.Unix(m.RevisionDate/1000, 0)

	rows, _ := db.pool.Query(ctx, `
		INSERT INTO summoner_records
			(summoner_id, account_id, puuid, revision_date, profile_icon_id, summoner_level)
		VALUES
			($1, $2, $3, $4, $5, $6)
		RETURNING *;
		`,
		m.Id, m.AccountId, m.Puuid, ts, m.ProfileIconId, m.SummonerLevel)

	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[SummonerRecord])
	if err != nil {
		return nil, err
	}

	return res, nil
}
