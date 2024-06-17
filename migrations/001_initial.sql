-- Write your migrate up statements here

create extension if not exists "uuid-ossp";

create table summonerrecords (
    record_id UUID default uuid_generate_v4() primary key,
    record_date TIMESTAMP default current_timestamp,

    account_id VARCHAR(56) not null,
    summoner_id VARCHAR(63) not null,
    puuid VARCHAR(78) not null,
    revision_date TIMESTAMP not null,
    summoner_level BIGINT not null,
    profile_icon_id INT not null
);

create table leaguerecords (
    record_id UUID default uuid_generate_v4() primary key,
    record_date TIMESTAMP default current_timestamp,

    summoner_id VARCHAR(63) not null,
    league_id VARCHAR(128) not null,
    tier VARCHAR(16) not null,
    division VARCHAR(8) not null,
    league_points INT not null,
    number_wins INT not null,
    number_losses INT not null
);

---- create above / drop below ----

drop table summonerrecords;
drop table leaguerecords;

drop extension if exists "uuid-ossp";

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
