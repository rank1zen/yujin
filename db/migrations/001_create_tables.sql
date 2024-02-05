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

CREATE INDEX ix_summoner_profile_name ON summoner_profile (name);
CREATE INDEX ix_summoner_records_puuid_ts ON summoner_records (name, puuid, record_date DESC);

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

---- create above / drop below ----

DROP TABLE summoner_profile;
DROP TABLE summoner_records;
DROP INDEX ix_summoner_profile_name;
DROP INDEX ix_summoner_records_puuid_ts;

DROP TABLE soloq_records;

DROP EXTENSION IF EXISTS "uuid-ossp";

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
