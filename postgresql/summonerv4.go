package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type SummonerRecordRow struct {
	RecordId      string    `db:"record_id"`
	RecordDate    time.Time `db:"record_date"`
	Puuid         string    `db:"puuid"`
	AccountId     string    `db:"account_id"`
	SummonerId    string    `db:"id"`
	Name          string    `db:"name"`
	ProfileIconId int       `db:"profile_icon_id"`
	SummonerLevel int       `db:"summoner_level"`
	RevisionDate  int       `db:"revision_date"`
}

type SummonerRecord struct {
	RecordId      string    `db:"record_id"`
	RecordDate    time.Time `db:"record_date"`
	Name          string    `db:"name"`
	ProfileIconId int       `db:"profile_icon_id"`
	SummonerLevel int       `db:"summoner_level"`
	RevisionDate  int       `db:"revision_date"`
}

type SummonerRecordArg struct {
	Puuid         string
	AccountId     string
	SummonerId    string
	Name          string
	ProfileIconId int
	SummonerLevel int
	RevisionDate  int
}

func (q *Queries) InsertSummonerRecord(ctx context.Context, r *SummonerRecordArg, date time.Time) (string, error) {
	const query = `
	INSERT INTO summoner_records
	(record_date, account_id, profile_icon_id, revision_date, name, id, puuid, summoner_level)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING record_id
	`

	var uuid string
	err := q.db.QueryRow(ctx, query,
		date,
		r.AccountId,
		r.ProfileIconId,
		r.RevisionDate,
		r.Name,
		r.SummonerId,
		r.Puuid,
		r.SummonerLevel,
	).Scan(&uuid)
	if err != nil {
		return "", fmt.Errorf("query error: %w", err)
	}

	return uuid, nil
}

func (q *Queries) DeleteSummonerRecord(ctx context.Context, id string) (string, error) {
	query := `
	DELETE FROM summoner_records
	WHERE record_id = $1
	RETURNING record_id
	`

	var uuid string
	err := q.db.QueryRow(ctx, query, id).Scan(&uuid)
	if err != nil {
		return "", fmt.Errorf("query error: %w", err)
	}

	return uuid, nil
}

func (q *Queries) SelectSummonerRecord(ctx context.Context, id string) (*SummonerRecord, error) {
	query := `
	SELECT record_id, record_date, name, profile_icon_id, summoner_level, revision_date
	FROM summoner_records
	WHERE record_id = $1
	`

	rows, err := q.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[SummonerRecord])
	if err != nil {
		return nil, fmt.Errorf("row error: %w", err)
	}

	return &record, nil
}

func (q *Queries) CountSummonerRecordsByPuuid(ctx context.Context, puuid string) (int64, error) {
	query := `
	SELECT count(*)
	FROM summoner_records
	WHERE puuid = $1
	`

	var count int64
	err := q.db.QueryRow(ctx, query, puuid).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (q *Queries) SelectSummonerRecordJoin(ctx context.Context, name string) {
	_ = `
	SELECT * FROM summoner_profile s
	INNER JOIN LATERAL (
	SELECT record_id, record_date, name, profile_icon_id, summoner_level, revision_date
	FROM summoner_records
	WHERE name = s.name
	ORDER BY record_date DESC
	LIMIT 1
	) l ON TRUE
	ORDER BY s.name DESC;
	`
}

func (q *Queries) SelectSummonerRecentByPuuid(ctx context.Context, puuid string) (*SummonerRecord, error) {
	query := `
	SELECT record_id, record_date, name, profile_icon_id, summoner_level, revision_date
	FROM summoner_records
	WHERE puuid = $1
	ORDER BY record_date DESC
	LIMIT 1
	`

	rows, err := q.db.Query(ctx, query, puuid)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[SummonerRecord])
	if err != nil {
		return nil, fmt.Errorf("row error: %w", err)
	}

	return &record, nil
}

func (q *Queries) SelectSummonerRecordsByPuuid(ctx context.Context, puuid string) (*[]SummonerRecord, error) {
	query := `
	SELECT record_id, record_date, name, profile_icon_id, summoner_level, revision_date
	FROM summoner_records
	WHERE puuid = $1
	ORDER BY record_date DESC
	LIMIT $2 OFFSET $3
	`

	rows, err := q.db.Query(ctx, query, puuid, 10, 0)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	records, err := pgx.CollectRows(rows, pgx.RowToStructByName[SummonerRecord])
	if err != nil {
		return nil, fmt.Errorf("row error: %w", err)
	}

	return &records, nil
}

func (q *Queries) SelectSummonerRecordsByName(ctx context.Context, name string) (*[]SummonerRecord, error) {
	query := `
	SELECT record_id, record_date, name, profile_icon_id, summoner_level, revision_date
	FROM summoner_records
	WHERE name = $1
	ORDER BY record_date DESC
	LIMIT $2 OFFSET $3
	`

	rows, err := q.db.Query(ctx, query, name, 10, 0)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	records, err := pgx.CollectRows(rows, pgx.RowToStructByName[SummonerRecord])
	if err != nil {
		return nil, fmt.Errorf("row error: %w", err)
	}

	return &records, nil
}
