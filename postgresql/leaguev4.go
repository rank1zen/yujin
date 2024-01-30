package postgresql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type SoloqRecord struct {
	RecordId   string
	RecordDate time.Time

	LeagueId     string
	SummonerName string
	Tier         string
	Rank         string
	League       string
	Wins         int32
	Losses       int32
}

func (q *Queries) InsertSoloqRecord(ctx context.Context, r *SoloqRecord) (string, error) {
	var id string

	query := `
	INSERT INTO soloq_records
	(record_date, league_id, summoner_name, tier, rank, league, wins, losses)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING record_id
	`

	err := q.db.QueryRow(ctx, query,
		r.RecordDate,
		r.LeagueId,
		r.SummonerName,
		r.Tier,
		r.Rank,
		r.League,
		r.Wins,
		r.Losses,
	).Scan(&id)

	return id, err
}

func (q *Queries) SelectSoloqRecordsBySummonerId(ctx context.Context, puuid string) ([]SoloqRecord, error) {
	query := `
	SELECT *
	FROM soloq_records
	WHERE summoner_id = $1
	ORDER BY record_date DESC
	LIMIT $2 OFFSET $3;
	`
	rows, _ := q.db.Query(ctx, query, puuid)
	defer rows.Close()

	records, err := pgx.CollectRows(rows, pgx.RowToStructByName[SoloqRecord])
	return records, err
}

func (q *Queries) SelectSoloqRecord(ctx context.Context, id string) (SoloqRecord, error) {
	query := `
	SELECT * 
	`

	rows, _ := q.db.Query(ctx, query, id)
	defer rows.Close()

	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[SoloqRecord])
	return record, err
}

func (q *Queries) CountSoloqRecordsBySummonerId(ctx context.Context, id string) (int64, error) {
	query := `
	SELECT COUNT(*)
	FROM soloq_records
	WHERE summoner_id = $1
	`

	var count int64
	err := q.db.QueryRow(ctx, query, id).Scan(&count)
	return count, err
}

func (q *Queries) DeletSoloqRecord(ctx context.Context, id string) (string, time.Time, error) {
	query := `
	DELETE FROM soloq_records
	WHERE record_id = $1
	RETURNING record_date, summoner_name
	`
	var s string
	var t time.Time
	err := q.db.QueryRow(ctx, query, id).Scan(&s, &t)
	return s, t, err
}

func (q *Queries) SelectSoloqRecordsByName(ctx context.Context, name string) ([]SoloqRecord, error) {
	query := `
	SELECT *
	FROM soloq_records
	WHERE summoner_name = $1
	ORDER BY record_date DESC
	LIMIT $2 OFFSET $3;
	`

	rows, _ := q.db.Query(ctx, query)
	defer rows.Close()

	records, err := pgx.CollectRows(rows, pgx.RowToStructByName[SoloqRecord])
	return records, err
}
