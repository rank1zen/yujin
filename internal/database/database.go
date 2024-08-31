package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/yujin/internal/ddragon"
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

	err = queryInsertRow(ctx, db.pool, "riot_accounts", vals)
	if err != nil {
		return nil, fmt.Errorf("inserting: %w", err)
	}

	return getAccount(gamename, tagline)
}

func (db *DB) CheckFirstTimeProfile(ctx context.Context, name string) (bool, error) {
	ids, err := db.GetAccount(ctx, name)
	if err != nil {
		return false, fmt.Errorf("getting account: %w", err)
	}

	var exists bool
	err = db.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM profile_summaries WHERE puuid = $1)`, ids.Puuid).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
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

	err = queryInsertRow(ctx, db.pool, "summoner_records", row)
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

	err = queryInsertRow(ctx, db.pool, "league_records", row)
	if err != nil {
		return err
	}

	return nil
}

type ProfileSummary struct {
	Name                string
	TagLine             string
	ProfileIconImageUrl string
	LastUpdated         string
	SoloqRank           string
	LeaguePoints        string
	Wins                string
	Losses              string
	SummonerLevel       int
}

func (db *DB) GetProfileSummary(ctx context.Context, name string) (*ProfileSummary, error) {
	ids, err := db.GetAccount(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("getting account: %w", err)
	}

	found, err := db.CheckFirstTimeProfile(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("checking profile: %w", err)
	}

	if !found {
		err := db.UpdateProfile(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("updating: %w", err)
		}
	}

	m := ProfileSummary{Name: ids.Name, TagLine: ids.TagLine}
	var iconId int
	// TODO: have seperated fields to wins and losses
	err = db.pool.QueryRow(ctx, `
	SELECT
		profile_icon_id,
		summoner_level,
		CASE WHEN 1=1
			AND tier          IS NOT NULL
			AND division      IS NOT NULL
			AND league_points IS NOT NULL
			AND number_wins   IS NOT NULL
			AND number_losses IS NOT NULL
		THEN FORMAT('%s %s %s LP, %s-%s', tier, division, league_points, number_wins, number_losses)
		ELSE 'Unranked'
		END AS soloq_rank
	FROM
		profile_summaries
	WHERE
		puuid = $1;
	`, ids.Puuid).Scan(&iconId, &m.SummonerLevel, &m.SoloqRank)
	if err != nil {
		return nil, fmt.Errorf("getting profile: %w", err)
	}

	m.ProfileIconImageUrl = ddragon.GetSummonerProfileIconUrl(iconId)

	return &m, nil
}

type ProfileMatch struct {
	PlayerWin    bool          `db:"player_win"`
	LpChange     int           `db:"-"`
	GameDate     time.Time     `db:"game_date"`
	GameDuration time.Duration `db:"game_duration"`

	Kills      string `db:"kills"`
	Deaths     string `db:"deaths"`
	Assists    string `db:"assists"`
	CreepScore string `db:"creep_score"`
	CsPer10    string `db:"cs_per_10"`
	Damage     string `db:"damage"`

	ItemIds         []int  `db:"items"`
	SpellIds        []int  `db:"spells"`
	ChampionName    string `db:"champion_name"`
	RunePrimaryId   int    `db:"rune_primary"`
	RuneSecondaryId int    `db:"rune_secondary"`
}

func (db *DB) GetProfileMatchList(ctx context.Context, name string, page int, ensure bool) (ProfileMatchList, error) {
	ids, err := db.GetAccount(ctx, name)
	if err != nil {
		return nil, err
	}

	start, count := 10*page, 10
	if ensure {
		err := db.ensureMatchlist(ctx, ids.Puuid, start, count)
		if err != nil {
			return nil, err
		}
	}

	rows, _ := db.pool.Query(ctx, `
	SELECT
		player_win,
		game_date,
		game_duration,
		kills,
		deaths,
		assists,
		creep_score,
		TO_CHAR(60 * creep_score / EXTRACT(epoch FROM game_duration), 'FM99999.0') AS cs_per_10,
		total_damage_dealt_to_champions AS damage,
		champion_name,
		array[item0_id, item1_id, item2_id, item3_id, item4_id, item5_id] as items,
		array[spell1_id, spell2_id] as spells,
		rune_primary_keystone AS rune_primary,
		rune_secondary_path AS rune_secondary
	FROM
		profile_matches
	WHERE
		puuid = $1
	OFFSET $2 LIMIT $3;
	`, ids.Puuid, start, count)

	matchlist, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[ProfileMatch])
	if err != nil {
		return nil, fmt.Errorf("getting matchlist: %w", err)
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
