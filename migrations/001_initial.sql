-- Write your migrate up statements here

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE SummonerRecords (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	record_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

	account_id      VARCHAR(56)  NOT NULL,
	summoner_id     VARCHAR(63)  NOT NULL,
	puuid           VARCHAR(78)  NOT NULL,
	revision_date   TIMESTAMP    NOT NULL,
	summoner_level  BIGINT       NOT NULL,
	profile_icon_id INT          NOT NULL
);

CREATE TABLE LeagueRecords (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	record_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

	summoner_id   VARCHAR(63)  NOT NULL,
	league_id     VARCHAR(128) NOT NULL,
	tier          VARCHAR(16)  NOT NULL,
	division      VARCHAR(8)   NOT NULL,
	league_points INT          NOT NULL,
	number_wins   INT          NOT NULL,
	number_losses INT          NOT NULL
);

---- create above / drop below ----

DROP TABLE SummonerRecords;
DROP TABLE LeagueRecords;

DROP EXTENSION IF EXISTS "uuid-ossp";

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
