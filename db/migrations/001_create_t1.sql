-- Write your migrate up statements here

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE summoners (
	id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	puuid TEXT NOT NULL,
	account_id TEXT NOT NULL,
	summoner_id TEXT NOT NULL,
	level BIGINT NOT NULL,
    profile_icon_id INT NOT NULL,
    name TEXT NOT NULL,
    last_revision TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    time_stamp TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

---- create above / drop below ----

DROP TABLE summoners;

DROP EXTENSION IF EXISTS "uuid-ossp";

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
