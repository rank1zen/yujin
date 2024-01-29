package postgresql

import (
	"context"
	"time"
)

const (
	insertSoloqRecord = `
	INSERT INTO soloq_records
	(record_date, league_id, summoner_name, tier, rank, league, wins, losses)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING record_id
	`

	countSoloqRecordsById = `
	SELECT COUNT(*)
	FROM soloq_records
	WHERE summoner_id = $1
	`

	deleteSoloqRecord = `
	DELETE FROM soloq_records
	WHERE record_id = $1
	RETURNING record_date, summoner_name
	`

	countSoloqRecordsBySummonerId = `
	SELECT *
	FROM soloq_records
	WHERE summoner_id = $1
	ORDER BY record_date DESC
	LIMIT $2 OFFSET $3;
	`

	coutnSoloqRecordsByName = `
	SELECT *
	FROM soloq_records
	WHERE summoner_name = $1
	ORDER BY record_date DESC
	LIMIT $2 OFFSET $3;
	`
)

type SoloqRecord struct {
	RecordId   string
	RecordDate time.Time

	LeagueId     string
	SummonerName string
	Tier         string
	Rank         string
	League       string
	wins         int32
	losses       int32
}

func InsertSoloqRecord(ctx context.Context, r *SoloqRecord) (string, error) {
	return "", nil
}

func CountSoloqRecordsById(ctx context.Context, id string) (int64, error) {
	return 0, nil
}

func DeletSoloqRecord(ctx context.Context, id string) (string, time.Time, error) {
	return "", time.Time{}, nil
}
