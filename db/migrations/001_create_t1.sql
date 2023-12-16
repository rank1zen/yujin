-- Write your migrate up statements here

CREATE TABLE summoners (
	id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
	puuid text NOT NULL,
	account_id text NOT NULL,
	summoner_id text NOT NULL,
	level BIGINT,
    profile_icon_id INT,
    name text NOT NULL,
    last_revision TIMESTAMP WITHOUT TIME ZONE,
    time_stamp TIMESTAMP WITHOUT TIME ZONE
);

---- create above / drop below ----

DROP TABLE summoners;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
