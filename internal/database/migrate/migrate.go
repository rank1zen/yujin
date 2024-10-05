package migrate

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func Migrate(ctx context.Context, conn *pgx.Conn) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for i, f := range migrations {
		err := f(tx)
		if err != nil {
			return fmt.Errorf("[migration %d]: %w", i, err)
		}
	}

	return tx.Commit(ctx)
}

var schemaVersion = len(migrations)

var migrations = []func(tx pgx.Tx) error{
	func(tx pgx.Tx) (err error) {
		sql := `
		CREATE TABLE schema_version (
			version text not null
		);

		CREATE DOMAIN riot_puuid       AS CHAR(78);
		CREATE DOMAIN riot_summoner_id AS VARCHAR(63);
		CREATE DOMAIN riot_account_id  AS VARCHAR(56);
		CREATE DOMAIN riot_match_id    AS VARCHAR(60);

		CREATE TABLE profiles (
			puuid        riot_puuid  primary key,
			name         varchar(32) not null,
			tagline      varchar(10) not null,
			last_updated timestamptz not null
		);

		CREATE TABLE summoner_records (
			record_id       uuid             default gen_random_uuid() primary key,
			record_date     timestamptz      default current_timestamp,
			account_id      riot_account_id  not null,
			summoner_id     riot_summoner_id not null,
			puuid           riot_puuid       not null,
			revision_date   timestamptz      not null,
			summoner_level  bigint           not null,
			profile_icon_id int              not null
		);

		CREATE TABLE league_records (
			record_id     uuid             default gen_random_uuid() primary key,
			record_date   timestamptz      default current_timestamp,
			summoner_id   riot_summoner_id not null,
			league_id     varchar(128),
			tier          varchar(16),
			division      varchar(8),
			league_points int,
			wins          int,
			losses        int
		);

		CREATE VIEW summoner_records_latest AS
		SELECT DISTINCT ON (puuid)
			summoner_id, puuid, summoner_level, profile_icon_id
		FROM summoner_records
		ORDER BY puuid, record_date DESC;

		CREATE VIEW league_records_latest AS
		SELECT DISTINCT ON (summoner_id)
			summoner_id, tier, division, league_points, wins, losses
		FROM league_records
		ORDER BY summoner_id, record_date DESC;

		CREATE VIEW profile_headers AS
		SELECT
			profile.name,
			profile.tagline,
			profile.last_updated,
			summoner.puuid,
			summoner.profile_icon_id,
			summoner.summoner_level,
			league.tier,
			league.division,
			league.league_points,
			league.wins,
			league.losses
		FROM summoner_records_latest AS summoner
		JOIN league_records_latest AS league ON summoner.summoner_id = league.summoner_id
		JOIN profiles AS profile ON summoner.puuid = profile.puuid;

		CREATE FUNCTION format_rank(tier varchar(16), div varchar(8), lp int) RETURNS varchar(32) AS $$
		BEGIN
			IF tier is not null AND div is not null AND lp is not null THEN
				CASE tier WHEN 'CHALLENGER', 'MASTER', 'GRANDMASTER' THEN
					RETURN format('%s %sLP', initcap(tier), lp);
				ELSE
					RETURN format('%s %s %sLP', initcap(tier), div, lp);
				END CASE;
			ELSE
				RETURN 'Unranked';
			END IF;
		END;
		$$ LANGUAGE plpgsql;

		CREATE FUNCTION format_win_loss(wins int, losses int) RETURNS varchar(32) AS $$
		BEGIN
			IF wins is not null AND losses is not null THEN
				RETURN format('%s-%s', wins, losses);
			ELSE
				RETURN '0-0';
			END IF;
		END;
		$$ LANGUAGE plpgsql;
		`
		_, err = tx.Exec(context.Background(), sql)
		return err
	},
	func(tx pgx.Tx) (err error) {
		sql := `
		CREATE TABLE matches (
			id            riot_match_id primary key,
			data_version  text          not null,
			game_date     timestamptz   not null,
			game_duration interval      not null,
			game_patch    varchar(32)   not null
		);

		CREATE TABLE match_participants (
			match_id riot_match_id not null,
			FOREIGN KEY(match_id)
				REFERENCES matches(id)
				ON DELETE CASCADE,
			team_id int         not null,
			id      int         not null,
			name    varchar(50) not null,
			puuid   riot_puuid  not null,
			unique(match_id, puuid),
			unique(match_id, id),

			player_position varchar(10) not null,
			champion_level  int         not null,
			champion_id     int         not null,
			champion_name   varchar(30) not null,
			kills           int         not null,
			deaths          int         not null,
			assists         int         not null,
			creep_score     int         not null,
			vision_score    int         not null,
			gold_earned     int         not null,
			gold_spent      int         not null,

			items                 int[7] not null,
			spell1_id             int    not null,
			spell2_id             int    not null,
			rune_primary_path     int    not null,
			rune_primary_keystone int    not null,
			rune_primary_slot1    int    not null,
			rune_primary_slot2    int    not null,
			rune_primary_slot3    int    not null,
			rune_secondary_path   int    not null,
			rune_secondary_slot1  int    not null,
			rune_secondary_slot2  int    not null,
			rune_shard_slot1      int    not null,
			rune_shard_slot2      int    not null,
			rune_shard_slot3      int    not null,

			physical_damage_dealt              int not null,
			physical_damage_dealt_to_champions int not null,
			physical_damage_taken              int not null,
			magic_damage_dealt                 int not null,
			magic_damage_dealt_to_champions    int not null,
			magic_damage_taken                 int not null,
			true_damage_dealt                  int not null,
			true_damage_dealt_to_champions     int not null,
			true_damage_taken                  int not null,
			total_damage_dealt                 int not null,
			total_damage_dealt_to_champions    int not null,
			total_damage_taken                 int not null
		);

		CREATE TABLE match_teams (
			match_id riot_match_id not null,
			FOREIGN KEY (match_id)
				REFERENCES matches(id)
				ON DELETE CASCADE,
			id int not null,
			unique (match_id, id),
			win boolean not null
		);

		CREATE VIEW profile_matches AS
		SELECT
			player.match_id,
			player.team_id,
			player.id,
			player.name,
			player.puuid,
			player.kills,
			player.deaths,
			player.assists,
			player.vision_score,
			player.creep_score,
			player.gold_earned,
			player.champion_level,
			player.champion_name,
			player.champion_id,
			player.total_damage_dealt_to_champions,
			player.spell1_id,
			player.spell2_id,
			player.items,
			player.rune_primary_keystone,
			player.rune_secondary_path,
			match.game_date,
			match.game_duration,
			match.game_patch,
			team.win
		FROM match_participants AS player
		JOIN matches AS match ON match.id = player.match_id
		JOIN match_teams AS team ON match.id = team.match_id AND player.team_id = team.id;
		`
		_, err = tx.Exec(context.Background(), sql)
		return err
	},
	func(tx pgx.Tx) (err error) {
		sql := `
		CREATE FUNCTION format_cs_per10(cs int, game_duration interval) RETURNS char(4) AS $$
		BEGIN
			RETURN TO_CHAR(60 * cs / EXTRACT(epoch FROM game_duration), 'FM99.0');
		END;
		$$ LANGUAGE plpgsql;

		CREATE FUNCTION format_kill_participation() RETURNS char(3) AS $$
		BEGIN
			RETURN '20%';
		END;
		$$ LANGUAGE plpgsql;

		CREATE FUNCTION format_damage_relative() RETURNS char(3) AS $$
		BEGIN
			RETURN '80%';
		END;
		$$ LANGUAGE plpgsql;

		CREATE FUNCTION get_champion_icon_url(id int) RETURNS varchar(128) AS $$
		BEGIN
			RETURN FORMAT('https://cdn.communitydragon.org/14.16.1/champion/%s/square', id);
		END;
		$$ LANGUAGE plpgsql;

		CREATE FUNCTION get_item_icon_urls(ids int[]) RETURNS text[] AS $$
		DECLARE
			urls text[] := array_fill(NULL::text, ARRAY[7]);
		BEGIN
			IF array_length(ids, 1) != 7 THEN
				RAISE EXCEPTION 'must have exactly 7 items';
			END IF;

			FOR i IN 1..7 LOOP
				IF ids[i] != 0 THEN
					urls[i] := FORMAT('https://ddragon.leagueoflegends.com/cdn/14.16.1/img/item/%s.png', ids[i]);
				END IF;
			END LOOP;
			RETURN urls;
		END;
		$$ LANGUAGE plpgsql;
		`
		_, err = tx.Exec(context.Background(), sql)
		return err
	},
}
