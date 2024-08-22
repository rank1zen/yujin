package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/yujin/pkg/riot"
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

	riot := riot.NewClient()

	return &DB{
		pool: pool,
		riot: riot,
	}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}

type Identifiers struct {
	Puuid      string
	SummonerId string
	Name       string
	TagLine    string
}

func (db *DB) GetAccount(ctx context.Context, name string) (*Identifiers, error) {
	parts := strings.SplitN(name, "-", 2)
	gamename, tagline := parts[0], parts[1]
	var found bool
	err := db.pool.QueryRow(ctx, `
	SELECT EXISTS (SELECT 1 FROM accounts WHERE name = $1 AND tagline = $2)
	`, gamename, tagline).Scan(&found)
	if err != nil {
		return nil, err
	}

	if found {
		var ids Identifiers
		err := db.pool.QueryRow(ctx, `
		SELECT
			summoner_id,
			puuid
		FROM
			accounts
		WHERE
			name = $1 AND tagline = $2;
		`, gamename, tagline).Scan(&ids.SummonerId, &ids.Puuid)
		if err != nil {
			return nil, err
		}

		return &ids, nil
	}

	acc, err := db.riot.GetAccountByRiotId(ctx, gamename, tagline)
	if err != nil {
		return nil, fmt.Errorf("fetching account: %w", err)
	}

	summ, err := db.riot.GetSummoner(ctx, acc.Puuid)
	if err != nil {
		return nil, fmt.Errorf("fetching summoner: %w", err)
	}

	vals := map[string]any{"puuid": summ.Puuid, "summoner_id": summ.Id}
	err = queryInsertRow(ctx, db.pool, "accounts", vals)
	if err != nil {
		return nil, fmt.Errorf("inserting: %w", err)
	}

	return &Identifiers{
		Puuid:      acc.Puuid,
		SummonerId: summ.Id,
	}, nil
}

func (db *DB) UpdateProfile(ctx context.Context, name string) error {
	ids, err := db.GetAccount(ctx, name)
	if err != nil {
		return fmt.Errorf("getting account: %w", err)
	}

	err = db.ensureMatchlist(ctx, ids.Puuid, 0, 10)
	if err != nil {
		return fmt.Errorf("matchlist: %w", err)
	}

	summoner, err := db.riot.GetSummoner(ctx, ids.Puuid)
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

	err = queryInsertRow(ctx, db.pool, "summoner_records", row)
	if err != nil {
		return err
	}

	leagues, err := db.riot.GetLeagueEntriesForSummoner(ctx, ids.SummonerId)
	if err != nil {
		return fmt.Errorf("getting league: %w", err)
	}

	soloq := findSoloqRank(leagues)
	row = map[string]any{"summoner_id": ids.SummonerId}
	if soloq != nil {
		row["league_id"] = soloq.LeagueId
		row["tier"] = soloq.Tier
		row["division"] = soloq.Rank
		row["league_points"] = soloq.LeaguePoints
		row["number_wins"] = soloq.Wins
		row["number_losses"] = soloq.Losses

	}

	err = queryInsertRow(ctx, db.pool, "league_records", row)
	if err != nil {
		return err
	}

	return nil
}

type ProfileSummary struct {
	ProfileIconId int32
	SummonerLevel int32
	Name          string
	Rank          string
}

func (db *DB) GetProfileSummary(ctx context.Context, name string) (map[string]any, error) {
	ids, err := db.GetAccount(ctx, name)
	if err != nil {
		return nil, err
	}

	rows, _ := db.pool.Query(ctx, `
	SELECT
		summoner.profile_icon_id,
		summoner.summoner_level,

		league.tier,
		league.division,
		league.league_points,
		league.number_wins,
		league.number_losses
	FROM
		summoner_records_newest AS summoner
	JOIN
		league_records_newest AS league
	ON
		summoner.summoner_id = league.summoner_id
	WHERE
		summoner.puuid = $1;
	`, ids.Puuid)

	// fn := func(row pgx.CollectableRow) (*ProfileSummary, error) {
	// 	var m ProfileSummary
	// 	var tier, div, lp, wins, losses *int
	// 	err := row.Scan(&m.ProfileIconId, &m.SummonerLevel, &tier, &div, &lp, &wins, &losses)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	if tier != nil {
	// 		m.Rank = fmt.Sprintf("%s %s %d LP %d-%d", *tier, *div, *lp, *wins, *losses)
	// 	} else {
	// 		m.Rank = "Unranked"
	// 	}
	//
	// 	return &m, nil
	// }

	m, err := pgx.CollectExactlyOneRow(rows, pgx.RowToMap)
	if err != nil {
		return nil, err
	}

	return m, nil
}

type ProfileMatch struct {
	GameDate     time.Time     `db:"game_date"`
	GameDuration time.Duration `db:"game_duration"`
	MatchId      string        `db:"match_id"`
	GamePatch    string        `db:"game_patch"`
	PlayerWin    bool          `db:"player_win"`
	PostGame     *MatchSummonerPostGame
}

func (db *DB) GetProfileMatchList(ctx context.Context, puuid string, page int, ensure bool) (ProfileMatchList, error) {
	start, count := 10*page, 10
	if ensure {
		err := db.ensureMatchlist(ctx, puuid, start, count)
		if err != nil {
			return nil, err
		}
	}

	rows, _ := db.pool.Query(ctx, `
	SELECT
		info.match_id,
		info.game_date,
		info.game_duration,
		info.game_patch,

		player.player_win,
		player.player_position,

		player.kills,
		player.deaths,
		player.assists,
		player.creep_score,
		player.gold_earned,
		player.gold_spent,
		player.champion_level,
		player.champion_id,
		player.champion_name,
		player.vision_score,

		player.item0_id,
		player.item1_id,
		player.item2_id,
		player.item3_id,
		player.item4_id,
		player.item5_id,
		player.item6_id,

		player.spell0_id,
		player.spell1_id,

		player.rune_primary_path,
		player.rune_primary_keystone,
		player.rune_primary_slot1,
		player.rune_primary_slot2,
		player.rune_primary_slot3,
		player.rune_secondary_path,
		player.rune_secondary_slot1,
		player.rune_secondary_slot2,
		player.rune_shard_slot1,
		player.rune_shard_slot2,
		player.rune_shard_slot3
	FROM
		match_info_records AS info
	JOIN
		match_participant_records AS player
	ON
		info.match_id = player.match_id
	WHERE
		puuid = $1
	OFFSET $2 LIMIT $3;
	`, puuid, start, count)

	fn := func(row pgx.CollectableRow) (*ProfileMatch, error) {
		var m ProfileMatch
		err := row.Scan(&m.MatchId, &m.GameDate, &m.GameDuration, &m.GamePatch)
		if err != nil {
			return nil, err
		}

		postgame, err := pgx.RowToAddrOfStructByName[MatchSummonerPostGame](row)
		if err != nil {
			return nil, err
		}

		m.PostGame = postgame
		return &m, err
	}

	matchlist, err := pgx.CollectRows(rows, fn)
	if err != nil {
		return nil, err
	}

	return matchlist, nil
}

// TODO
// func (db *DB) GetProfileMatchSummary(ctx context.Context, name RiotName, matchID RiotMatchId) (*ProfileMatchSummary, error) {
// 	ids, err := db.GetAccount(ctx, name)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var m ProfileMatchSummary
// 	rows, _ := db.pool.Query(ctx, `
// 	SELECT
// 		player_position,
// 		kills,
// 		deaths,
// 		assists,
// 		creep_score,
// 		champion_level,
// 		champion_id,
// 		vision_score,
// 		items_arr,
// 		spells_arr
// 	FROM
// 		match_participant_simple
// 	WHERE
// 		match_id = $1;
// 	`, matchID)
// 	players, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[MatchSummonerPostGame])
// 	if err != nil {
// 		return nil, fmt.Errorf("select post games: %w", err)
// 	}
//
// 	if len(players) != 10 {
// 		return nil, fmt.Errorf("got players: %v", len(players))
// 	}
//
// 	rows, _ = db.pool.Query(ctx, `
// 	SELECT
// 		rune_slot,
// 		rune_id
// 	FROM
// 		match_rune_records
// 	WHERE
// 		puuid = $1 AND match_id = $2
// 	`, ids.Puuid, matchID)
//
// 	return &m, nil
// }
