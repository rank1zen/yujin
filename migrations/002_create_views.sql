-- Write your migrate up statements here

CREATE VIEW summoner_records_newest AS
WITH numbered_records AS (
    SELECT *, row_number() OVER (PARTITION BY puuid ORDER BY record_date DESC) AS rn
    FROM summoner_records
)
SELECT
    record_id, record_date, account_id, summoner_id, puuid, revision_date, summoner_level, profile_icon_id
FROM numbered_records
WHERE rn = 1;

CREATE VIEW league_records_newest AS
WITH numbered_records AS (
    SELECT *, row_number() OVER (PARTITION BY summoner_id ORDER BY record_date) AS rn
    FROM league_records
)
SELECT
    record_id, record_date, summoner_id, league_id, tier, division, league_points, number_wins, number_losses
FROM numbered_records
WHERE rn = 1;


CREATE VIEW summoner_profile AS
SELECT
    s.summoner_level, s.profile_icon_id,
    l.tier, l.division, l.league_points, l.number_wins, l.number_losses
FROM summoner_records_newest AS s
INNER JOIN league_records_newest AS l ON s.summoner_id = l.summoner_id;


CREATE VIEW match_participant_simple AS
WITH
items_agg AS (
    SELECT puuid, array_agg(item_id ORDER BY item_slot) AS items_arr
    FROM match_item_records GROUP BY puuid
),
spells_agg AS (
    SELECT puuid, array_agg(spell_id ORDER BY spell_slot) AS spells_arr
    FROM match_summonerspell_records GROUP BY puuid
)
SELECT
    info.match_id, info.game_date, info.game_duration, info.game_patch,
    p.puuid, p.player_win, p.player_position, p.kills, p.deaths, p.assists,
    p.creep_score, p.champion_level, p.champion_id, p.vision_score,
    p.participant_id,
    items.items_arr, spells.spells_arr,
    runes.keystone, runes.secondary
FROM match_info_records AS info
INNER JOIN match_participant_records AS p ON info.match_id = p.match_id
INNER JOIN items_agg AS items ON p.puuid = items.puuid
INNER JOIN spells_agg AS spells ON p.puuid = spells.puuid
INNER JOIN rune_simple AS runes ON p.puuid = runes.puuid;

CREATE VIEW runes_simple AS
SELECT
    match_id,
    puuid,
    MAX(CASE WHEN rune_type = 'main keystone' THEN rune_id) AS keystone,
    MAX(CASE WHEN rune_type = 'secondary path' THEN rune_id) AS secondary
FROM
    match_rune_records
GROUP BY match_id, puuid;

CREATE VIEW runes_full AS
SELECT
    match_id,
    puuid,
    MAX(CASE WHEN rune_type = 'main keystone' THEN rune_id) AS main_keystone,
    MAX(CASE WHEN rune_type = 'main slot1' THEN rune_id) AS main_slot1,
    MAX(CASE WHEN rune_type = 'main slot2' THEN rune_id) AS main_slot2,
    MAX(CASE WHEN rune_type = 'main slot3' THEN rune_id) AS main_slot3,
    MAX(CASE WHEN rune_type = 'secondary slot1' THEN rune_id) AS secondary_slot1,
    MAX(CASE WHEN rune_type = 'secondary slot2' THEN rune_id) AS secondary_slot2,
    MAX(CASE WHEN rune_type = 'shard slot1' THEN rune_id) AS shard_slot1,
    MAX(CASE WHEN rune_type = 'shard slot2' THEN rune_id) AS shard_slot2,
    MAX(CASE WHEN rune_type = 'shard slot3' THEN rune_id) AS shard_slot3
FROM match_rune_records GROUP BY match_id;

---- create above / drop below ----

DROP VIEW summoner_records_newest;
DROP VIEW summoner_profile;
DROP VIEW league_records_newest;
DROP VIEW match_participant_simple;
DROP VIEW runes_full;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
