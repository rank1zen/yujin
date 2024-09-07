package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/yujin/internal/database/migrate"
	"github.com/rank1zen/yujin/internal/pgxutil"
	"github.com/rank1zen/yujin/internal/riot"
	"go.opentelemetry.io/otel/trace"
)

type DB struct {
	pool   *pgxpool.Pool
	riot   *riot.Client
	tracer trace.Tracer
}

func NewDB(ctx context.Context, url string) (*DB, error) {
	pgxCfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres connection string: %w", err)
	}

	pgxCfg.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		return conn.Ping(ctx) == nil
	}

	pgxCfg.ConnConfig.Tracer = &tracer{}

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("connecting to postgres: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = migrate.Migrate(ctx, conn.Conn())
	if err != nil {
		return nil, err
	}

	riot := &riot.Client{}

	return &DB{
		pool: pool,
		riot: riot,
	}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}

type Account struct {
	Puuid      string
	SummonerId string
	Name       string
	TagLine    string
}

func (db *DB) GetAccount(ctx context.Context, name string) (*Account, error) {
	getAccount := func(gamename, tagline string) (*Account, error) {
		var ids Account
		err := db.pool.QueryRow(ctx, `
		SELECT
			summoner_id,
			puuid,
			name,
			tagline
		FROM
			riot_accounts
		WHERE 1=1
			AND UPPER(name) = UPPER($1)
			AND UPPER(tagline) = UPPER($2);
		`, gamename, tagline).Scan(&ids.SummonerId, &ids.Puuid, &ids.Name, &ids.TagLine)
		if err != nil {
			return nil, fmt.Errorf("getting db: %w", err)
		}
		return &ids, nil
	}

	parts := strings.SplitN(name, "-", 2)
	gamename, tagline := parts[0], parts[1]
	var found bool
	err := db.pool.QueryRow(ctx, `
	SELECT EXISTS (
		SELECT 1
		FROM
			riot_accounts
		WHERE 1=1
			AND UPPER(name) = UPPER($1)
			AND UPPER(tagline) = UPPER($2)
	);
	`, gamename, tagline).Scan(&found)
	if err != nil {
		return nil, fmt.Errorf("checking db: %w", err)
	}

	if found {
		return getAccount(gamename, tagline)
	}

	acc, err := db.riot.GetAccountByRiotId(ctx, gamename, tagline)
	if err != nil {
		return nil, fmt.Errorf("fetching account: %w", err)
	}

	summ, err := db.riot.GetSummonerByPuuid(ctx, acc.Puuid)
	if err != nil {
		return nil, fmt.Errorf("fetching summoner: %w", err)
	}

	vals := map[string]any{
		"puuid":       summ.Puuid,
		"summoner_id": summ.Id,
		"name":        acc.GameName,
		"tagline":     acc.TagLine,
	}

	err = pgxutil.QueryInsertRow(ctx, db.pool, "riot_accounts", vals)
	if err != nil {
		return nil, fmt.Errorf("inserting: %w", err)
	}

	return getAccount(gamename, tagline)
}

func (db *DB) UpdateProfile(ctx context.Context, name string) error {
	ids, err := db.GetAccount(ctx, name)
	if err != nil {
		return fmt.Errorf("getting account: %w", err)
	}

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	summoner, err := db.riot.GetSummonerByPuuid(ctx, ids.Puuid)
	if err != nil {
		return fmt.Errorf("getting summoner: %w", err)
	}

	row := map[string]any{
		"summoner_id":     summoner.Id,
		"account_id":      summoner.AccountId,
		"puuid":           summoner.Puuid,
		"revision_date":   time.Unix(summoner.RevisionDate/1000, 0), // NOTE: might want to double check this
		"profile_icon_id": summoner.ProfileIconId,
		"summoner_level":  summoner.SummonerLevel,
	}

	err = pgxutil.QueryInsertRow(ctx, tx, "summoner_records", row)
	if err != nil {
		return err
	}

	leagues, err := db.riot.GetLeagueEntriesForSummoner(ctx, ids.SummonerId)
	if err != nil {
		return fmt.Errorf("getting league: %w", err)
	}

	row = map[string]any{"summoner_id": ids.SummonerId}
	for _, entry := range leagues {
		if entry.QueueType == "RANKED_SOLO_5x5" {
			row["league_id"] = entry.LeagueId
			row["tier"] = entry.Tier
			row["division"] = entry.Rank
			row["league_points"] = entry.LeaguePoints
			row["number_wins"] = entry.Wins
			row["number_losses"] = entry.Losses
		}
	}

	err = pgxutil.QueryInsertRow(ctx, tx, "league_records", row)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

type ProfileSummary struct {
	Name          string `db:"name"`
	TagLine       string `db:"tagline"`
	LastUpdated   string `db:"last_updated"`
	SoloqRank     string `db:"soloq_rank"`
	WinLoss       string `db:"win_loss"`
	SummonerLevel int    `db:"summoner_level"`
}

func (db *DB) GetProfileSummary(ctx context.Context, name string) (*ProfileSummary, error) {
	ids, err := db.GetAccount(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("getting account: %w", err)
	}

	var exists bool
	err = db.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM profile_summaries WHERE puuid = $1)`, ids.Puuid).Scan(&exists)
	if err != nil {
		return nil, err
	}

	if !exists {
		err := db.UpdateProfile(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("updating: %w", err)
		}
	}

	rows, _ := db.pool.Query(ctx, `
	SELECT
		name,
		tagline,
		TO_CHAR(record_date, 'YYYY MM-DD HH24:MI') AS last_updated,
		summoner_level,
		CASE
			WHEN 1=1 AND tier IS NOT NULL AND division IS NOT NULL AND league_points IS NOT NULL THEN
			CASE WHEN tier IN ('CHALLENGER', 'MASTER', 'GRANDMASTER') THEN
				FORMAT('%s %s LP', INITCAP(tier), league_points)
			ELSE FORMAT('%s %s %s LP', INITCAP(tier), division, league_points)
			END
		ELSE 'Unranked'
		END AS soloq_rank,
		CASE WHEN 1=1
			AND number_wins   IS NOT NULL
			AND number_losses IS NOT NULL
		THEN FORMAT('%s-%s', number_wins, number_losses)
		ELSE '0-0'
		END AS win_loss
	FROM
		profile_summaries
	WHERE
		puuid = $1;
	`, ids.Puuid)

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[ProfileSummary])
}

type ProfileMatch struct {
	PlayerWin    bool   `db:"win"`
	LpChange     int    `db:"-"`
	GameDate     string `db:"game_date"`
	GameDuration string `db:"game_duration"`

	Kills      string `db:"kills"`
	Deaths     string `db:"deaths"`
	Assists    string `db:"assists"`
	CreepScore string `db:"creep_score"`
	CsPer10    string `db:"cs_per_10"`
	Damage     string `db:"damage"`

	ItemIds         []int  `db:"items"`
	SpellIds        []int  `db:"spells"`
	ChampionIconUrl string `db:"champion_icon_url"`
	RunePrimaryId   int    `db:"rune_primary"`
	RuneSecondaryId int    `db:"rune_secondary"`
}

type ProfileMatchList []*ProfileMatch

func (db *DB) GetProfileMatchList(ctx context.Context, name string, page int, ensure bool) (ProfileMatchList, error) {
	ids, err := db.GetAccount(ctx, name)
	if err != nil {
		return nil, err
	}

	start, count := 10*page, 10
	if ensure {
		err := ensureMatchList(ctx, db.pool, db.riot, ids.Puuid, start, count)
		if err != nil {
			return nil, fmt.Errorf("ensuring matchlist: %w", err)
		}
	}

	rows, _ := db.pool.Query(ctx, `
	SELECT
		win,
		TO_CHAR(game_date, 'MM-DD HH24:MI') AS game_date,
		EXTRACT(MINUTE FROM game_duration) || 'm ' || EXTRACT(SECOND FROM game_duration) || 's' AS game_duration,
		kills,
		deaths,
		assists,
		creep_score,
		TO_CHAR(60 * creep_score / EXTRACT(epoch FROM game_duration), 'FM99999.0') AS cs_per_10,
		total_damage_dealt_to_champions AS damage,
		FORMAT('https://cdn.communitydragon.org/14.16.1/champion/%s/square', champion_id) AS champion_icon_url,
		array[item0_id, item1_id, item2_id, item3_id, item4_id, item5_id] as items,
		array[spell1_id, spell2_id] as spells,
		rune_primary_keystone AS rune_primary,
		rune_secondary_path AS rune_secondary
	FROM
		profile_matches
	WHERE
		puuid = $1
	ORDER BY
		game_date DESC
	OFFSET $2 LIMIT $3;
	`, ids.Puuid, start, count)

	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[ProfileMatch])
}
