-- Write your migrate up statements here

CREATE TABLE MatchInfoRecords (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	record_date TIMESTAMP NOT NULL,
	match_id VARCHAR(64),
	game_date TIMESTAMP NOT NULL,
	game_duration INTERVAL NOT NULL,
	game_patch VARCHAR(128),
  UNIQUE (match_id)
);

CREATE TABLE MatchTeamRecords (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	match_id VARCHAR(64),
	team_id INT NOT NULL,
  team_win BOOLEAN NOT NULL,
  team_surrendered BOOLEAN NOT NULL,
  team_early_surrendered BOOLEAN NOT NULL,
  FOREIGN KEY (match_id) REFERENCES MatchInfoRecords (match_id),
	UNIQUE (match_id, team_id)
);

CREATE TABLE MatchBanRecords (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	match_id VARCHAR(64),
  team_id INT NOT NULL,
	champion_id INT NOT NULL,
	turn INT NOT NULL,
  FOREIGN KEY (match_id) REFERENCES MatchInfoRecords (match_id),
  FOREIGN KEY (match_id, team_id) REFERENCES MatchTeamRecords (match_id, team_id)
);

CREATE TABLE MatchObjectiveRecords (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	match_id VARCHAR(64) NOT NULL,
  team_id INT NOT NULL,
	name VARCHAR(128) NOT NULL,
	first BOOLEAN NOT NULL,
	kills INT NOT NULL,
  FOREIGN KEY (match_id) REFERENCES MatchInfoRecords (match_id),
  FOREIGN KEY (match_id, team_id) REFERENCES MatchTeamRecords (match_id, team_id)
);

CREATE TABLE MatchParticipantRecords (
  record_id UUID default uuid_generate_v4() PRIMARY KEY,
  match_id VARCHAR(64),
  puuid VARCHAR(128),
  player_win BOOLEAN,
  player_position VARCHAR(16),
  kills INT,
  deaths INT,
  assists INT,
  creep_score INT,
  gold_earned INT,
  champion_level INT,
  champion_id INT,
  FOREIGN KEY (match_id) REFERENCES MatchInfoRecords (match_id),
  UNIQUE (match_id, puuid)
);

CREATE TABLE MatchRuneRecords (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  match_id VARCHAR(64) REFERENCES MatchInfoRecords (match_id),
  puuid VARCHAR(128),
  FOREIGN KEY (match_id, puuid) REFERENCES MatchParticipantRecords (match_id, puuid)
);

CREATE TABLE MatchSummonerSpellRecords (
	record_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  match_id VARCHAR(64) NOT NULL REFERENCES MatchInfoRecords(match_id),
  puuid VARCHAR(128),
  spell_slot INT,
  spell_id INT,
  FOREIGN KEY (match_id, puuid) REFERENCES MatchParticipantRecords (match_id, puuid)
);

CREATE TABLE MatchItemRecords (
  record_id UUID default uuid_generate_v4() PRIMARY KEY,
  match_id VARCHAR(64) not null REFERENCES MatchInfoRecords (match_id),
  puuid VARCHAR(128),
  item_id INT,
  item_slot INT,
  FOREIGN KEY (match_id, puuid) REFERENCES MatchParticipantRecords (match_id, puuid)
);

---- create above / drop below ----

DROP TABLE MatchInfoRecords;
DROP TABLE MatchTeamRecords;
DROP TABLE MatchBanRecords;
DROP TABLE MatchObjectiveRecords;
DROP TABLE MatchParticipantRecords;
DROP TABLE MatchRuneRecords;
DROP TABLE MatchSummonerSpellRecords;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
