-- name: SelectRecentRecordForSummoner :one
SELECT
    level,
    profile_icon_id,
    name,
    last_revision,
    time_stamp
FROM summoners
WHERE puuid = $1
ORDER BY time_stamp DESC
LIMIT 1;

-- name: SelectRecordsForSummoner :many
SELECT
    level,
    profile_icon_id,
    name,
    last_revision,
    time_stamp
FROM summoners
WHERE puuid = $1
ORDER BY time_stamp DESC
LIMIT $2 OFFSET $3;

-- name: CountSummonerRecords :one
SELECT count(*) FROM summoners
WHERE puuid = $1;

-- name: InsertSummoner :one
INSERT INTO summoners (puuid, account_id, summoner_id, level, profile_icon_id, name, last_revision, time_stamp)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;

-- name: DeleteSummoner :exec
DELETE FROM summoners WHERE id = $1;
