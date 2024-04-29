-- Write your migrate up statements here

ALTER TABLE MatchParticipantRecords DROP COLUMN win;
ALTER TABLE MatchParticipantRecords ADD COLUMN summoner_name VARCHAR(32);
ALTER TABLE MatchParticipantRecords ADD COLUMN team_id INT;
ALTER TABLE MatchParticipantRecords RENAME COLUMN id to participant_id;
ALTER TABLE MatchParticipantRecords ALTER COLUMN position type VARCHAR(16);
ALTER TABLE MatchParticipantRecords ALTER COLUMN champion_name type VARCHAR(32);

---- create above / drop below ----

ALTER TABLE MatchParticipantRecords DROP COLUMN summoner_name;
ALTER TABLE MatchParticipantRecords DROP COLUMN team_id;
ALTER TABLE MatchParticipantRecords RENAME participant_id to id;
ALTER TABLE MatchParticipantRecords ALTER COLUMN position type TEXT;
ALTER TABLE MatchParticipantRecords ALTER COLUMN champion_name type TEXT;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
