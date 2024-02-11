package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LeagueV4Query struct {
	db *pgxpool.Pool
}

type SoloqRecord struct {
	RecordId   string    `db:"record_id"`
	RecordDate time.Time `db:"record_date"`
	Name       string    `db:"summoner_name"`
	Tier       string    `db:"tier"`
	Rank       string    `db:"rank"`
	Lp         int32     `db:"league_points"`
	Wins       int32     `db:"wins"`
	Losses     int32     `db:"losses"`
}

type SoloqRecordArg struct {
	LeagueId   string
	SummonerId string
	Name       string
	Tier       string
	Rank       string
	Lp         int
	Wins       int
	Losses     int
}

func (q *LeagueV4Query) InsertSoloqRecord(ctx context.Context, r *SoloqRecordArg, ts time.Time) (string, error) {
	query := `
	INSERT INTO soloq_records
	(record_date, league_id, summoner_id, summoner_name, tier, rank, league_points, wins, losses)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING record_id
	`

	var id string
	err := q.db.QueryRow(ctx, query,
		ts,
		r.LeagueId,
		r.SummonerId,
		r.Name,
		r.Tier,
		r.Rank,
		r.Lp,
		r.Wins,
		r.Losses,
	).Scan(&id)
	return id, err
}

func (q *LeagueV4Query) SelectSoloqRecord(ctx context.Context, id string) (*SoloqRecord, error) {
	query := `
	SELECT record_id, record_date, summoner_name, tier, rank, league_points, wins, losses
	FROM soloq_records
	WHERE record_id = $1
	`

	rows, _ := q.db.Query(ctx, query, id)
	defer rows.Close()

	record, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[SoloqRecord])
	return &record, err
}

func (q *LeagueV4Query) DeleteSoloqRecord(ctx context.Context, id string) error {
	query := `
	DELETE FROM soloq_records
	WHERE record_id = $1
	`

	_, err := q.db.Exec(ctx, query, id)
	return err
}

func (q *LeagueV4Query) SelectSoloqRecordsByName(ctx context.Context, name string) (*[]SoloqRecord, error) {
	query := `
	SELECT record_id, record_date, summoner_name, tier, rank, league_points, wins, losses
	FROM soloq_records
	WHERE summoner_name = $1
	ORDER BY record_date DESC
	LIMIT $2 OFFSET $3
	`
	rows, _ := q.db.Query(ctx, query, name, 10, 0)
	defer rows.Close()

	records, err := pgx.CollectRows(rows, pgx.RowToStructByName[SoloqRecord])
	return &records, err
}

func (q *LeagueV4Query) SelectSoloqRecordNewestByName(ctx context.Context, name string) (*SoloqRecord, error) {
	query := `
	SELECT record_id, record_date, summoner_name, tier, rank, league_points, wins, losses
	FROM soloq_records
	WHERE summoner_name = $1
	ORDER BY record_date DESC
	LIMIT 1
	`

	rows, _ := q.db.Query(ctx, query)
	defer rows.Close()

	records, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[SoloqRecord])
	return &records, err
}

func (q *LeagueV4Query) CountSoloqRecordsBySummonerId(ctx context.Context, id string) (int64, error) {
	query := `
	SELECT COUNT(*)
	FROM soloq_records
	WHERE summoner_id = $1
	`

	var count int64
	err := q.db.QueryRow(ctx, query, id).Scan(&count)
	return count, err
}
