-- name: InsertSoloqRecord :one
INSERT INTO soloq_records
(
    record_date,
    league_id,
    summoner_id,
    summoner_name,
    tier,
    rank,
    league_points,
    wins,
    losses
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING record_id;

-- name: DeleteSoloqRecord :one
DELETE FROM soloq_records
WHERE record_id = $1
RETURNING record_date, summoner_name;

-- name: SelectSoloqRecordsBySummonerId :many
SELECT *
FROM soloq_records
WHERE summoner_id = $1
ORDER BY record_date DESC
LIMIT $2 OFFSET $3;

-- name: SelectSoloqRecordsByName :many
SELECT *
FROM soloq_records
WHERE summoner_name = $1
ORDER BY record_date DESC
LIMIT $2 OFFSET $3;

-- name: CountSoloqRecordsById :one
SELECT COUNT(*)
FROM soloq_records
WHERE summoner_id = $1;
