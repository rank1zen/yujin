-- Write your migrate up statements here

CREATE TABLE MatchInfoRecords (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	record_date TIMESTAMP NOT NULL,

	match_id VARCHAR(64) UNIQUE NOT NULL,
	game_date TIMESTAMP NOT NULL,
	game_duration INTERVAL,
	game_patch VARCHAR(128)
);

CREATE TABLE MatchTeamRecords (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	match_id VARCHAR(64) NOT NULL REFERENCES MatchRecords(match_id),
	team_id INT NOT NULL,
	UNIQUE (match_id, team_id),

    team_win BOOLEAN,
    team_surrendered BOOLEAN,
    team_early_surrendered BOOLEAN,
);

CREATE TABLE MatchBanRecords (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	match_id VARCHAR(64) NOT NULL REFERENCES MatchRecords(match_id),
    team_id INT NOT NULL REFERENCES MatchTeamRecords(team_id),

	champion_id INT,
	turn INT
);

CREATE TABLE MatchObjectiveRecords (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	match_id VARCHAR(64) NOT NULL REFERENCES MatchRecords (match_id),
    team_id INT NOT NULL REFERENCES MatchTeamRecords(team_id),

	name VARCHAR(128),
	first BOOLEAN,
	kills INT
);

CREATE TABLE MatchParticipantRecords (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	match_id VARCHAR(64) NOT NULL REFERENCES MatchRecords (match_id),
	puuid VARCHAR(128) NOT NULL,
	UNIQUE (match_id, puuid),

	player_win BOOLEAN,
	player_position VARCHAR(16),

	kills INT,
	deaths INT,
	assists INT,
	creep_score INT,
	gold_earned INT,

	champion_level      INT,
	champion_id         INT
);

CREATE TABLE MatchRuneRecords (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,

    match_id VARCHAR(64) NOT NULL,
    puuid VARCHAR(128) NOT NULL
);

CREATE TABLE MatchSummonerSpellRecords (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,

    match_id VARCHAR(64) NOT NULL,
    puuid VARCHAR(128) NOT NULL,
    spell_id INT
);

---- create above / drop below ----

DROP TABLE MatchRecords;
DROP TABLE MatchTeamRecords;
DROP TABLE MatchBanRecords;
DROP TABLE MatchObjectiveRecords;
DROP TABLE MatchParticipantRecords;
DROP TABLE MatchRuneRecords;
DROP TABLE MatchSummonerSpellRecords;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
