-- Write your migrate up statements here

CREATE DOMAIN riot_puuid AS CHAR(78);
CREATE DOMAIN riot_summoner_id AS VARCHAR(63);
CREATE DOMAIN riot_account_id AS VARCHAR(56);
CREATE DOMAIN riot_match_id AS VARCHAR(60);

CREATE TABLE summoner_records (
    record_id UUID default gen_random_uuid() primary key,
    record_date TIMESTAMP default current_timestamp,
    account_id riot_account_id NOT NULL,
    summoner_id riot_summoner_id not null,
    puuid riot_puuid not null,
    revision_date TIMESTAMP not null,
    summoner_level BIGINT not null,
    profile_icon_id INT not null
);

CREATE VIEW summoner_records_newest AS
WITH numbered_records AS (
    SELECT *, row_number() OVER (PARTITION BY puuid ORDER BY record_date DESC) AS rn
    FROM summoner_records
)
SELECT
    record_id,
    record_date,
    account_id,
    summoner_id,
    puuid,
    revision_date,
    summoner_level,
    profile_icon_id
FROM
    numbered_records
WHERE
    rn = 1;

CREATE TABLE league_records (
    record_id UUID default gen_random_uuid() primary key,
    record_date TIMESTAMP default current_timestamp,
    summoner_id riot_summoner_id not null,
    league_id VARCHAR(128),
    tier VARCHAR(16),
    division VARCHAR(8),
    league_points INT,
    number_wins INT,
    number_losses INT
);

CREATE VIEW league_records_newest AS
WITH numbered_records AS (
    SELECT *, row_number() OVER (PARTITION BY summoner_id ORDER BY record_date DESC) AS rn
    FROM league_records
)
SELECT
    record_id, record_date, summoner_id, league_id, tier, division, league_points, number_wins, number_losses
FROM numbered_records
WHERE rn = 1;

CREATE TABLE match_info_records (
	record_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
	match_id riot_match_id NOT NULL,
	game_date TIMESTAMP NOT NULL,
	game_duration INTERVAL NOT NULL,
	game_patch VARCHAR(128) NOT NULL,
    UNIQUE (match_id)
);

CREATE TABLE match_team_records (
	record_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
	match_id riot_match_id NOT NULL,
	team_id INT NOT NULL,
    team_win BOOLEAN NOT NULL,
    team_surrendered BOOLEAN NOT NULL,
    team_early_surrendered BOOLEAN NOT NULL,
    FOREIGN KEY (match_id) REFERENCES match_info_records (match_id),
	UNIQUE (match_id, team_id)
);

CREATE TABLE match_ban_records (
	record_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
	match_id riot_match_id NOT NULL,
    team_id INT NOT NULL,
	champion_id INT NOT NULL,
	turn INT NOT NULL,
    FOREIGN KEY (match_id) REFERENCES match_info_records (match_id),
    FOREIGN KEY (match_id, team_id) REFERENCES match_team_records (match_id, team_id)
);

CREATE TABLE match_objective_records (
	record_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    match_id riot_match_id NOT NULL,
    team_id INT NOT NULL,
    name VARCHAR(128) NOT NULL,
	first BOOLEAN NOT NULL,
    kills INT NOT NULL,
    FOREIGN KEY (match_id) REFERENCES match_info_records (match_id),
    FOREIGN KEY (match_id, team_id) REFERENCES match_team_records (match_id, team_id)
);

CREATE TABLE match_participant_records (
    record_id UUID default gen_random_uuid() PRIMARY KEY,
    match_id riot_match_id NOT NULL,
    puuid VARCHAR(78) NOT NULL,
    team_id INT NOT NULL,
    player_win BOOLEAN NOT NULL,
    participant_id INT,
    player_position VARCHAR(10) NOT NULL,
    kills INT NOT NULL,
    deaths INT NOT NULL,
    assists INT NOT NULL,
    creep_score INT NOT NULL,
    vision_score INT NOT NULL,
    gold_earned INT NOT NULL,
    gold_spent INT,
    double_kills INT,
    triple_kills INT,
    quadra_kills INT,
    penta_kills INT,
    champion_level INT NOT NULL,
    champion_id INT NOT NULL,
    detector_wards_placed INT,
    sight_wards_bought_ingame INT,
    vision_wards_bought_ingame INT,
    wards_killed INT,
    wards_placed INT,
    total_damage_dealt_to_champions INT,

    rune_main_path INT NOT NULL,
    rune_main_keystone INT NOT NULL,
    rune_main_slot1 INT NOT NULL,
    rune_main_slot2 INT NOT NULL,
    rune_main_slot3 INT NOT NULL,
    rune_secondary_path INT NOT NULL,
    rune_secondary_slot1 INT NOT NULL,
    rune_secondary_slot2 INT NOT NULL,
    rune_shard_slot1 INT NOT NULL,
    rune_shard_slot2 INT NOT NULL,
    rune_shard_slot3 INT NOT NULL,

    FOREIGN KEY (match_id) REFERENCES match_info_records (match_id),
    UNIQUE (match_id, puuid)
);

CREATE TABLE match_summonerspell_records (
	record_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    match_id riot_match_id NOT NULL,
    puuid riot_puuid NOT NULL,
    spell_casts INT,
    spell_slot INT,
    spell_id INT,
    FOREIGN KEY (match_id) REFERENCES match_info_records (match_id),
    FOREIGN KEY (match_id, puuid) REFERENCES match_participant_records (match_id, puuid)
);

CREATE TABLE match_item_records (
    record_id UUID default gen_random_uuid() PRIMARY KEY,
    match_id riot_match_id NOT NULL,
    puuid riot_puuid NOT NULL,
    item_id INT,
    item_slot INT,
    FOREIGN KEY (match_id) REFERENCES match_info_records (match_id),
    FOREIGN KEY (match_id, puuid) REFERENCES match_participant_records (match_id, puuid)
);

CREATE VIEW match_summoner_postgame AS
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
    info.match_id,
    info.game_date,
    info.game_duration,
    info.game_patch,
    p.puuid,
    p.player_win,
    p.player_position,
    p.kills,
    p.deaths,
    p.assists,
    p.creep_score,
    p.champion_level,
    p.champion_id,
    p.vision_score,
    p.participant_id,
    p.rune_main_keystone,
    p.rune_secondary_path,
    items.items_arr,
    spells.spells_arr
FROM match_info_records AS info
INNER JOIN match_participant_records AS p ON info.match_id = p.match_id
INNER JOIN items_agg AS items ON p.puuid = items.puuid
INNER JOIN spells_agg AS spells ON p.puuid = spells.puuid;

CREATE TABLE match_champion_kill_event_records (
    record_id UUID default gen_random_uuid() PRIMARY KEY,
    match_id riot_match_id NOT NULL,
	timestamp INTERVAL NOT NULL,
    pos_x INT NOT NULL,
    pos_y INT NOT NULL,
    bounty INT NOT NULL,
    shutdown_bounty INT NOT NULL,
    killer_id INT NOT NULL,
    victim_id INT NOT NULL,
    FOREIGN KEY (match_id) REFERENCES match_info_records (match_id)
);

CREATE TABLE match_item_event_records (
    record_id UUID default gen_random_uuid() PRIMARY KEY,
    match_id riot_match_id NOT NULL,
	timestamp INTERVAL NOT NULL,
    participant_id INT NOT NULL,
    item_id INT NOT NULL,
    FOREIGN KEY (match_id) REFERENCES match_info_records (match_id)
);

CREATE TABLE match_spell_event_records (
    record_id UUID default gen_random_uuid() PRIMARY KEY,
    match_id riot_match_id NOT NULL,
	timestamp INTERVAL NOT NULL,
    participant_id INT NOT NULL,
    skill_slot INT NOT NULL,
    FOREIGN KEY (match_id) REFERENCES match_info_records (match_id)
);

---- create above / drop below ----

DROP TABLE summoner_records;
DROP TABLE league_records;
DROP TABLE match_info_records;
DROP TABLE match_team_records;
DROP TABLE match_ban_records;
DROP TABLE match_objective_records;
DROP TABLE match_participant_records;
DROP TABLE match_rune_records;
DROP TABLE match_summonerspell_records;

DROP VIEW summoner_records_newest;
DROP VIEW league_records_newest;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
