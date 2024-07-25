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
    puuid VARCHAR(78),
    team_id INT,

    player_win BOOLEAN,
    participant_id INT,
    player_position VARCHAR(10),
    kills INT,
    deaths INT,
    assists INT,
    creep_score INT,
    vision_score INT,
    gold_earned INT,
    gold_spent INT,
    double_kills INT,
    triple_kills INT,
    quadra_kills INT,
    penta_kills INT,
    champion_level INT,
    champion_id INT,
    detector_wards_placed INT,
    sight_wards_bought_ingame INT,
    vision_wards_bought_ingame INT,
    wards_killed INT,
    wards_placed INT,
    total_damage_dealt_to_champions INT,

    FOREIGN KEY (match_id) REFERENCES match_info_records (match_id),
    UNIQUE (match_id, puuid)
);

CREATE TYPE rune_slot_type AS ENUM (
    'main path', 'main keystone', 'main slot1', 'main slot2', 'main slot3',
    'secondary path', 'secondary slot1', 'secondary slot2',
    'shard slot1', 'shard slot2', 'shard slot3',
);

CREATE TABLE match_rune_records (
	record_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    match_id riot_match_id NOT NULL,
    puuid riot_puuid NOT NULL,

    rune_slot rune_slot_type NOT NULL,
    rune_id int NOT NULL,

    FOREIGN KEY (match_id) REFERENCES match_info_records (match_id),
    FOREIGN KEY (match_id, puuid) REFERENCES match_participant_records (match_id, puuid),
    UNIQUE(match_id, puuid, rune_slot)
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

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
