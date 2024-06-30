-- Write your migrate up statements here

CREATE TABLE summoner_records (
  record_id UUID default gen_random_uuid() primary key,
  record_date TIMESTAMP default current_timestamp,
  account_id VARCHAR(56) not null,
  summoner_id VARCHAR(63) not null,
  puuid VARCHAR(78) not null,
  revision_date TIMESTAMP not null,
  summoner_level BIGINT not null,
  profile_icon_id INT not null
);

CREATE TABLE league_records (
  record_id UUID default gen_random_uuid() primary key,
  record_date TIMESTAMP default current_timestamp,
  summoner_id VARCHAR(63) not null,
  league_id VARCHAR(128),
  tier VARCHAR(16),
  division VARCHAR(8),
  league_points INT,
  number_wins INT,
  number_losses INT
);

CREATE TABLE match_info_records (
	record_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
	match_id VARCHAR(64),
	game_date TIMESTAMP NOT NULL,
	game_duration INTERVAL NOT NULL,
	game_patch VARCHAR(128),
  UNIQUE (match_id)
);

CREATE TABLE match_team_records (
	record_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
	match_id VARCHAR(64),
	team_id INT NOT NULL,
  team_win BOOLEAN NOT NULL,
  team_surrendered BOOLEAN NOT NULL,
  team_early_surrendered BOOLEAN NOT NULL,
  FOREIGN KEY (match_id) REFERENCES match_info_records (match_id),
	UNIQUE (match_id, team_id)
);

CREATE TABLE match_ban_records (
	record_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
	match_id VARCHAR(64),
  team_id INT NOT NULL,
	champion_id INT NOT NULL,
	turn INT NOT NULL,
  FOREIGN KEY (match_id) REFERENCES match_info_records (match_id),
  FOREIGN KEY (match_id, team_id) REFERENCES match_team_records (match_id, team_id)
);

CREATE TABLE match_objective_records (
	record_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
	match_id VARCHAR(64) NOT NULL,
  team_id INT NOT NULL,
	name VARCHAR(128) NOT NULL,
	first BOOLEAN NOT NULL,
	kills INT NOT NULL,
  FOREIGN KEY (match_id) REFERENCES match_info_records (match_id),
  FOREIGN KEY (match_id, team_id) REFERENCES match_team_records (match_id, team_id)
);

CREATE TABLE match_participant_records (
  record_id UUID default gen_random_uuid() PRIMARY KEY,
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
  FOREIGN KEY (match_id) REFERENCES match_info_records (match_id),
  UNIQUE (match_id, puuid)
);

CREATE TABLE match_rune_records (
	record_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  match_id VARCHAR(64) REFERENCES match_info_records (match_id),
  puuid VARCHAR(128),
  temp_name VARCHAR(16),
  rune_id int,
  FOREIGN KEY (match_id, puuid) REFERENCES match_participant_records (match_id, puuid)
);

CREATE TABLE match_summonerspell_records (
	record_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  match_id VARCHAR(64) NOT NULL REFERENCES match_info_records(match_id),
  puuid VARCHAR(128),
  spell_casts INT,
  spell_slot INT,
  spell_id INT,
  FOREIGN KEY (match_id, puuid) REFERENCES match_participant_records (match_id, puuid)
);

CREATE TABLE match_item_records (
  record_id UUID default gen_random_uuid() PRIMARY KEY,
  match_id VARCHAR(64) not null REFERENCES match_info_records (match_id),
  puuid VARCHAR(128),
  item_id INT,
  item_slot INT,
  FOREIGN KEY (match_id, puuid) REFERENCES match_participant_records (match_id, puuid)
);

---- create above / drop below ----

DROP TABLE summoner_records;
DROP TABLE league_records;
DROP TABLE match_info_records;
DROP TABLE match_team_records;
DROP TABLE match_ban_records;
DROP TABLE match_objective_records;
DROP TABLE match_participant_records;
DROP TABLE match_rune_records;
DROP TABLE match_summonerspell_records;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
