package migrate

import (
	"context"
	"embed"

	"github.com/jackc/pgx/v5"
)

//go:embed migrations/*.sql
var embeddedQueries embed.FS

func Migrate(ctx context.Context, conn *pgx.Conn) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, m := range migrations {
		err := m(tx)
		if err != nil {
			return err
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
			record_id uuid default gen_random_uuid() primary key,
			record_date timestamptz default current_timestamp,
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
}
