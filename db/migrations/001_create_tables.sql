-- Write your migrate up statements here

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE summoner_records (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	record_date TIMESTAMP,

	account_id TEXT,
	profile_icon_id INT,
	revision_date BIGINT,
	name TEXT,
	id TEXT, 
	puuid TEXT,
	summoner_level BIGINT
);

CREATE TABLE soloq_records (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	record_date TIMESTAMP,

	league_id TEXT,
	summoner_id TEXT, 
	summoner_name TEXT, 
	tier TEXT, 
	rank TEXT, 
	league_points INT,
	wins INT,
	losses INT
);

---- create above / drop below ----

DROP TABLE summoner_records;

DROP TABLE soloq_records;

DROP EXTENSION IF EXISTS "uuid-ossp";

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
