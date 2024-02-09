-- Write your migrate up statements here

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE summoner_profile (
    name TEXT PRIMARY KEY,
    puuid TEXT NOT NULL UNIQUE,
    account_id TEXT NOT NULL,
    summoner_id TEXT NOT NULL
);

CREATE TABLE summoner_records (
    record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    record_date TIMESTAMPTZ NOT NULL,

    account_id VARCHAR(150) NOT NULL,
    profile_icon_id INT NOT NULL,
    revision_date BIGINT NOT NULL,
    name VARCHAR(150) NOT NULL,
    id VARCHAR(150) NOT NULL,
    puuid VARCHAR(150) NOT NULL,
    summoner_level BIGINT NOT NULL
);


CREATE TABLE soloq_records (
    record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    record_date TIMESTAMPTZ NOT NULL,

    league_id TEXT NOT NULL,
    summoner_id TEXT NOT NULL,
    summoner_name TEXT NOT NULL,
    tier TEXT NOT NULL,
    rank TEXT NOT NULL,
    league_points INT NOT NULL,
    wins INT NOT NULL,
    losses INT NOT NULL
);

CREATE TYPE perk_style_selection AS (
    perk INT,
    var1 INT,
    var2 INT,
    var3 INT
);

CREATE TYPE perk_style AS (
    description TEXT,
    style INT,
    selections perk_style_selection []
);

CREATE TYPE perk_stat AS (
    defense INT,
    flex INT,
    offense INT
);

CREATE TYPE perk AS (
    --styles perk_style[],
    stat_perks perk_stat
);

CREATE TABLE match_v5 (
    match_id VARCHAR(100) PRIMARY KEY,
    runes PERK
);

CREATE TABLE participant (
	id INT,
	puuid VARCHAR(150),
	summoner_level BIGINT,
	win BOOL,
	position TEXT,
	kills INT,
	deaths INT,
	assists INT,
	creep_score INT,
	gold_earned INT,
	champion_id INT,
	champion_name TEXT,
	champion_level INT,
	item0 INT,
	item1 INT,
	item2 INT,
	item3 INT,
	item4 INT,
	item5 INT,
	item6 INT,
	vision_score INT,
	wards_placed INT,
	control_wards_placed INT,
	first_blood_assist BOOL,
	first_tower_assist BOOL,
	turret_takedowns INT,
	physical_damage_dealt_to_champions INT,
	magic_damage_dealt_to_champions INT,
	true_damage_dealt_to_champions INT,
	total_damage_dealt_to_champions INT,
	total_damage_taken INT,
	total_heals_on_teammates INT
);

CREATE INDEX ix_summoner_profile_name ON summoner_profile (name);
CREATE INDEX ix_summoner_records_puuid_ts ON summoner_records (

    name, puuid, record_date DESC
);

---- create above / drop below ----

DROP TABLE summoner_profile;
DROP TABLE summoner_records;
DROP INDEX ix_summoner_profile_name;
DROP INDEX ix_summoner_records_puuid_ts;

DROP TABLE soloq_records;

DROP EXTENSION IF EXISTS "uuid-ossp";

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
