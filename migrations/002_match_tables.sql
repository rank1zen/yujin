-- Write your migrate up statements here

CREATE TABLE MatchRecords (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	record_date TIMESTAMP NOT NULL,
	match_id VARCHAR(64) UNIQUE NOT NULL,

	start_ts TIMESTAMP,
	duration INTERVAL,
	surrender BOOLEAN,
	patch VARCHAR(128)
);

CREATE TABLE MatchTeamRecords (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	match_id VARCHAR(64) NOT NULL REFERENCES MatchRecords(match_id),
	team_id INT NOT NULL
);

CREATE TABLE MatchBanRecords (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	match_id VARCHAR(64) NOT NULL REFERENCES MatchRecords (match_id),
        team_id INT NOT NULL,

	champion_id INT,
	turn INT
);

CREATE TABLE MatchObjectiveRecords (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	match_id VARCHAR(64) NOT NULL REFERENCES MatchRecords (match_id),
        team_id INT NOT NULL,

	name VARCHAR(128),
	first BOOLEAN,
	kills INT
);

CREATE TABLE MatchParticipantRecords (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	match_id VARCHAR(64) NOT NULL REFERENCES MatchRecords (match_id),
	puuid VARCHAR(128),
	UNIQUE (match_id, puuid),

	id INT,
	summoner_level BIGINT,
	win BOOLEAN,
	position TEXT,
	kills INT,
	deaths INT,
	assists INT,
	creep_score INT,
	gold_earned INT,
	champion_id INT,
	champion_name TEXT,
	champion_level INT,
	vision_score INT,
	wards_placed INT,
	control_wards_placed INT,
	first_blood_assist BOOLEAN,
	first_tower_assist BOOLEAN,
	turret_takedowns INT,
	physical_damage_dealt_to_champions INT,
	magic_damage_dealt_to_champions INT,
	true_damage_dealt_to_champions INT,
	total_damage_dealt_to_champions INT,
	total_damage_taken INT,
	total_heals_on_teammates INT
);

---- create above / drop below ----

DROP TABLE MatchRecords;
DROP TABLE MatchTeamRecords;
DROP TABLE MatchBanRecords;
DROP TABLE MatchObjectiveRecords;
DROP TABLE MatchParticipantRecords;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
