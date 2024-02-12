-- Write your migrate up statements here

CREATE TYPE team_champion_ban AS (
	champion_id INT,
	turn INT
);

CREATE TYPE team_objective AS (
	name VARCHAR(128),
	first BOOLEAN,
	kills INT
);

CREATE TYPE participant_perk_style_selection AS (
	perk INT,
	var1 INT,
	var2 INT,
	var3 INT
);

CREATE TYPE participant_perk_style AS (
	description VARCHAR(128),
	style INT,
	selections participant_perk_style_selection[]
);

CREATE TYPE participant_perk_stat AS (
	defense INT,
	flex INT,
	offense INT
);

CREATE TYPE participant_perk AS (
	styles participant_perk_style[],
	stats participant_perk_stat
);

CREATE TABLE match_record (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	record_date TIMESTAMP NOT NULL,
	match_id VARCHAR(128) UNIQUE NOT NULL,

	start_ts TIMESTAMP NOT NULL,
	duration INTERVAL NOT NULL,
	surrender BOOLEAN NOT NULL,
	patch VARCHAR(128) NOT NULL
);

CREATE TABLE match_team_record (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	match_id VARCHAR(128) NOT NULL REFERENCES match_record(match_id),

	team_id INT NOT NULL,
	objectives team_objective[] NOT NULL,
	bans team_champion_ban[] NOT NULL
);

CREATE TABLE match_participant (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	match_id VARCHAR(128) NOT NULL REFERENCES match_record(match_id),
	puuid VARCHAR(128),
	UNIQUE (match_id, puuid),

	runes participant_perk,
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

DROP TABLE match_participant;
DROP TABLE match_team;
DROP TABLE match;

DROP TYPE participant_runes;
DROP TYPE participant_perk_style;
DROP TYPE participant_perk_style_selection;
DROP TYPE team_champion_ban;
DROP TYPE team_objective;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
