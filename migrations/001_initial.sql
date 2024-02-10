-- Write your migrate up statements here

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE summoner_records (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	record_date TIMESTAMP NOT NULL,

	account_id VARCHAR(128) NOT NULL,
	profile_icon_id INT NOT NULL,
	revision_date BIGINT NOT NULL,
	name VARCHAR(128) NOT NULL,
	id VARCHAR(128) NOT NULL,
	puuid VARCHAR(128) NOT NULL,
	summoner_level BIGINT NOT NULL
);

CREATE TABLE soloq_records (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	record_date TIMESTAMP NOT NULL,

	league_id VARCHAR(128) NOT NULL,
	summoner_id VARCHAR(128) NOT NULL,
	summoner_name VARCHAR(128) NOT NULL,
	tier VARCHAR(64) NOT NULL,
	division VARCHAR(64) NOT NULL,
	league_points INT NOT NULL,
	number_wins INT NOT NULL,
	number_losses INT NOT NULL
);

CREATE INDEX ix_summoner_records_puuid ON summoner_records (puuid, record_date DESC);
CREATE INDEX ix_soloq_records_summonerid ON soloq_records (summoner_id, record_date DESC);

---- create above / drop below ----

DROP TABLE summoner_records;
DROP TABLE soloq_records;

DROP INDEX ix_summoner_records_puuid;
DROP INDEX ix_soloq_records_summonerid;

DROP EXTENSION IF EXISTS "uuid-ossp";

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
