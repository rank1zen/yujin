package postgresql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type SummonerRecord struct {
	RecordDate    time.Time
	AccountId     string
	ProfileIconId int32
	RevisionDate  int64
	Name          string
	SummonerId    string
	Puuid         string
	SummonerLevel int64
}

const (
	insertSummonerRecord = `
	INSERT INTO summoner_records
	(record_date, account_id, profile_icon_id, revision_date, name, id, puuid, summoner_level)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING record_id
	`

	selectSummonerRecordById = `
	SELECT *
	FROM summoner_records
	WHERE id = $1
	`

	deleteSummonerRecord = `
	DELETE FROM summoner_records
	WHERE record_id = $1
	RETURNING record_date, name
	`

	countSummonerRecordsByPuuid = `
	SELECT COUNT(*)
	FROM summoner_records
	WHERE puuid = $1
	`

	selectSummonerRecordsByPuuid = `
	SELECT *
	FROM summoner_records
	WHERE puuid = $1
	ORDER BY record_date DESC
	LIMIT $2 OFFSET $3
	`
)

func (q *Queries) InsertSummonerRecord(ctx context.Context, rec *SummonerRecord) (string, error) {
	var id string
	err := q.db.QueryRow(ctx, insertSummonerRecord,
		rec.RecordDate,
		rec.AccountId,
		rec.ProfileIconId,
		rec.RevisionDate,
		rec.Name,
		rec.SummonerId,
		rec.Puuid,
		rec.SummonerLevel,
	).Scan(&id)

	return id, err
}

func (q *Queries) SelectSummonerRecordByPuuid(ctx context.Context, puuid string) ([]SummonerRecord, error) {
	rows, _ := q.db.Query(ctx, selectSummonerRecordsByPuuid, puuid)
	defer rows.Close()

	records, err := pgx.CollectRows(rows, pgx.RowToStructByName[SummonerRecord])
	if err != nil {
		return []SummonerRecord{}, err
	}

	return records, nil
}

func (q *Queries) SelectSummonerRecordById(ctx context.Context, id string) (SummonerRecord, error) {
	rows, _ := q.db.Query(ctx, selectSummonerRecordById, id)
	defer rows.Close()

	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[SummonerRecord])
	return record, err
}

func (q *Queries) CountSummonerRecordsByPuuid(ctx context.Context, puuid string) (int64, error) {
	var count int64
	err := q.db.QueryRow(ctx, countSummonerRecordsByPuuid, puuid).Scan(&count)
	return count, err
}

func (q *Queries) DeleteSummonerRecord(ctx context.Context, id string) (string, time.Time, error) {
	var s string
	var t time.Time
	err := q.db.QueryRow(ctx, deleteSummonerRecord, id).Scan(&s, &t)
	return s, t, err
}
