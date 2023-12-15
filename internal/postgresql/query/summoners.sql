-- name: ListAll :one
SELECT * FROM summoners
WHERE id = $1 LIMIT 1;

-- name: CreateSummoner :one
INSERT INTO summoners (
	puuid, account_id, summoner_id, level, profile_icon_id, name, last_revision, time_stamp
)
