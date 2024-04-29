-- Write your migrate up statements here

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE SummonerRecords (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	record_date TIMESTAMP NOT NULL,

	account_id VARCHAR(128),
	profile_icon_id INT,
	revision_date BIGINT,
	name VARCHAR(128),
	id VARCHAR(128),
	puuid VARCHAR(128),
	summoner_level BIGINT
);

CREATE TABLE LeagueRecords (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	record_date TIMESTAMP NOT NULL,

	league_id VARCHAR(128),
	summoner_id VARCHAR(128),
	summoner_name VARCHAR(128),
	tier VARCHAR(16),
	division VARCHAR(8),
	league_points INT,
	number_wins INT,
	number_losses INT
);

---- create above / drop below ----

DROP TABLE SummonerRecords;
DROP TABLE LeagueRecords;

DROP EXTENSION IF EXISTS "uuid-ossp";

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
