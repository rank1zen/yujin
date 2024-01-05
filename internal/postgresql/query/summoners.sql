-- name: SelectSummonerRecentByPuuid :one
SELECT *
FROM summoner_records
WHERE puuid = $1
ORDER BY record_date DESC
LIMIT 1;

-- name: SelectSummonerRecords :many
SELECT *
FROM summoner_records
ORDER BY record_date DESC
LIMIT $1
OFFSET $2;

-- name: SelectSummonerRecordsByPuuid :many
SELECT *
FROM summoner_records
WHERE puuid = $1
ORDER BY record_date DESC
LIMIT $2
OFFSET $3;

-- name: SelectSummonerRecordsNoIds :many
SELECT record_date, revision_date, name, summoner_level
FROM summoner_records
ORDER BY record_date DESC
LIMIT $1
OFFSET $2;

-- name: CountSummonerRecords :one
SELECT count(*)
FROM summoner_records;

-- name: CountSummonerRecordsByPuuid :one
SELECT count(*)
FROM summoner_records
WHERE puuid = $1;

-- name: InsertSummoner :one
INSERT INTO summoner_records (puuid, account_id, id, summoner_level, profile_icon_id, name, revision_date, record_date)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING record_id;

-- name: DeleteSummoner :exec
DELETE FROM summoner_records
WHERE record_id = $1;
