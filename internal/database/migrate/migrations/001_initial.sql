-- Write your migrate up statements here

CREATE DOMAIN riot_puuid       AS CHAR(78);
CREATE DOMAIN riot_summoner_id AS VARCHAR(63);
CREATE DOMAIN riot_account_id  AS VARCHAR(56);
CREATE DOMAIN riot_match_id    AS VARCHAR(60);

CREATE TABLE riot_accounts (
    record_id UUID default gen_random_uuid() primary key,

    summoner_id riot_summoner_id NOT NULL,
    puuid       riot_puuid       NOT NULL,
    name        VARCHAR(32)      NOT NULL,
    tagline     VARCHAR(32)      NOT NULL
);

CREATE TABLE summoner_records (
    record_id UUID default gen_random_uuid() primary key,
    record_date TIMESTAMPTZ default current_timestamp,

    account_id      riot_account_id  NOT NULL,
    summoner_id     riot_summoner_id NOT NULL,
    puuid           riot_puuid       NOT NULL,
    revision_date   TIMESTAMPTZ        NOT NULL,
    summoner_level  BIGINT           NOT NULL,
    profile_icon_id INT              NOT NULL
);

-- TODO: this query sucks
CREATE VIEW summoner_records_newest AS
WITH numbered_records AS (
    SELECT *, row_number() OVER (PARTITION BY puuid ORDER BY record_date DESC) AS rn
    FROM summoner_records
)
SELECT
    record_id, record_date, account_id, summoner_id, puuid, revision_date, summoner_level, profile_icon_id
FROM numbered_records
WHERE rn = 1;

CREATE TABLE league_records (
    record_id UUID default gen_random_uuid() primary key,
    record_date TIMESTAMPTZ default current_timestamp,

    summoner_id   riot_summoner_id NOT NULL,
    league_id     VARCHAR(128),
    tier          VARCHAR(16),
    division      VARCHAR(8),
    league_points INT,
    number_wins   INT,
    number_losses INT
);

-- TODO: this query sucks
CREATE VIEW league_records_newest AS
WITH numbered_records AS (
    SELECT *, row_number() OVER (PARTITION BY summoner_id ORDER BY record_date DESC) AS rn
    FROM league_records
)
SELECT
    record_id, record_date, summoner_id, league_id, tier, division, league_points, number_wins, number_losses
FROM numbered_records
WHERE rn = 1;

CREATE VIEW profile_summaries AS
SELECT
    account.name,
    account.tagline,

    summoner.puuid,
    summoner.profile_icon_id,
    summoner.summoner_level,

    league.record_date,
    league.tier,
    league.division,
    league.league_points,
    league.number_wins,
    league.number_losses
FROM
    summoner_records_newest AS summoner
JOIN
    league_records_newest AS league
ON
    summoner.summoner_id = league.summoner_id
JOIN
    riot_accounts AS account
ON
    summoner.puuid = account.puuid;

CREATE TABLE match_info_records (
	record_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,

	match_id      riot_match_id NOT NULL,
    data_version  TEXT          NOT NULL,
	game_date     TIMESTAMPTZ     NOT NULL,
	game_duration INTERVAL      NOT NULL,
	game_patch    VARCHAR(128)  NOT NULL,

    UNIQUE (match_id)
);

CREATE TABLE match_participant_records (
    record_id UUID default gen_random_uuid() PRIMARY KEY,

    match_id        riot_match_id NOT NULL,
    puuid           riot_puuid    NOT NULL,
    team_id         INT           NOT NULL,
    participant_id  INT           NOT NULL,

    player_position VARCHAR(10) NOT NULL,
    champion_level  INT         NOT NULL,
    champion_id     INT         NOT NULL,
    champion_name   VARCHAR(30) NOT NULL,

    kills        INT NOT NULL,
    deaths       INT NOT NULL,
    assists      INT NOT NULL,
    creep_score  INT NOT NULL,
    vision_score INT NOT NULL,

    spell1_id int NOT NULL,
    spell2_id int NOT NULL,

    item0_id int NOT NULL,
    item1_id int NOT NULL,
    item2_id int NOT NULL,
    item3_id int NOT NULL,
    item4_id int NOT NULL,
    item5_id int NOT NULL,
    item6_id int NOT NULL,

    rune_primary_path     INT NOT NULL,
    rune_primary_keystone INT NOT NULL,
    rune_primary_slot1    INT NOT NULL,
    rune_primary_slot2    INT NOT NULL,
    rune_primary_slot3    INT NOT NULL,
    rune_secondary_path   INT NOT NULL,
    rune_secondary_slot1  INT NOT NULL,
    rune_secondary_slot2  INT NOT NULL,
    rune_shard_slot1      INT NOT NULL,
    rune_shard_slot2      INT NOT NULL,
    rune_shard_slot3      INT NOT NULL,

    physical_damage_dealt              INT NOT NULL,
    physical_damage_dealt_to_champions INT NOT NULL,
    physical_damage_taken              INT NOT NULL,
    magic_damage_dealt                 INT NOT NULL,
    magic_damage_dealt_to_champions    INT NOT NULL,
    magic_damage_taken                 INT NOT NULL,
    true_damage_dealt                  INT NOT NULL,
    true_damage_dealt_to_champions     INT NOT NULL,
    true_damage_taken                  INT NOT NULL,
    total_damage_dealt                 INT NOT NULL,
    total_damage_dealt_to_champions    INT NOT NULL,
    total_damage_taken                 INT NOT NULL,

    FOREIGN KEY(match_id)
        REFERENCES match_info_records (match_id),
    UNIQUE (match_id, puuid)
);

CREATE TABLE match_team_records (
	record_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,

	match_id riot_match_id NOT NULL,
	team_id  INT           NOT NULL,
    win      BOOLEAN       NOT NULL,

    FOREIGN KEY (match_id) REFERENCES match_info_records (match_id),
	UNIQUE (match_id, team_id)
);

CREATE TABLE match_ban_records (
	record_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,

	match_id    riot_match_id NOT NULL,
    team_id     INT           NOT NULL,
	champion_id INT           NOT NULL,
	turn        INT           NOT NULL,

    FOREIGN KEY (match_id) REFERENCES match_info_records (match_id),
    FOREIGN KEY (match_id, team_id) REFERENCES match_team_records (match_id, team_id)
);

CREATE TABLE match_objective_records (
	record_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,

    match_id riot_match_id NOT NULL,
    team_id  INT           NOT NULL,
    name     VARCHAR(64)   NOT NULL,
	first    BOOLEAN       NOT NULL,
    kills    INT           NOT NULL,

    FOREIGN KEY (match_id) REFERENCES match_info_records (match_id),
    FOREIGN KEY (match_id, team_id) REFERENCES match_team_records (match_id, team_id)
);

CREATE VIEW profile_matches AS
SELECT
    info.match_id,
    info.game_date,
    info.game_duration,
    info.game_patch,

    player.puuid,
    player.kills,
    player.deaths,
    player.assists,
    player.creep_score,
    player.champion_level,
    player.champion_name,
    player.champion_id,
    player.total_damage_dealt_to_champions,
    player.spell1_id,
    player.spell2_id,
    player.item0_id,
    player.item1_id,
    player.item2_id,
    player.item3_id,
    player.item4_id,
    player.item5_id,
    player.rune_primary_keystone,
    player.rune_secondary_path,

    team.win
FROM
    match_info_records AS info
JOIN
    match_participant_records AS player
ON
    info.match_id = player.match_id
JOIN
    match_team_records AS team
ON 1=1
    AND info.match_id = team.match_id
    AND player.team_id = team.team_id;

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

DROP VIEW summoner_records_newest;
DROP VIEW league_records_newest;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
