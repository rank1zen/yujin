CREATE TABLE summoners (
	id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
	puuid VARCHAR,
	account_id VARCHAR,
	summoner_id VARCHAR,
	level BITINT,
    profile_icon_id INT,
    name VARCHAR,
    last_revision TIMESTAMP WITHOUT TIME ZONE,
    time_stamp TIMESTAMP WITHOUT TIME ZONE
);
