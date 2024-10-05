package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/internal/ddragon"
	"github.com/rank1zen/yujin/internal/logging"
	"github.com/rank1zen/yujin/internal/pgxutil"
	"github.com/rank1zen/yujin/internal/riot"
)

// returns whether it exists, will fetch from riot if it exists but not in db yet
func (db *DB) ProfileExists(ctx context.Context, puuid string) (bool, error) {
	var exists bool
	err := db.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM profiles WHERE puuid = $1)`, puuid).Scan(&exists)
	if err != nil {
		return false, err
	}

	if !exists {
		_, err := db.riot.AccountGetByPuuid(ctx, puuid)
		if err == nil {
			err := db.ProfileUpdate(ctx, puuid)
			if err != nil {
				return false, err
			}
			return true, nil
		} else if errors.Is(err, riot.ErrNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}

	return exists, nil
}

func (db *DB) ProfileUpdate(ctx context.Context, puuid string) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// remember to update profile table
	account, err := db.riot.AccountGetByPuuid(ctx, puuid)
	if err != nil {
		return err
	}

	_, err = db.pool.Exec(ctx, `
	INSERT INTO profiles (puuid, name, tagline, last_updated)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT(puuid)
	DO UPDATE SET name = $2, tagline = $3, last_updated = $4;
	`, puuid, account.GameName, account.TagLine, time.Now())
	if err != nil {
		return err
	}

	summoner, err := db.riot.GetSummonerByPuuid(ctx, puuid)
	if err != nil {
		return fmt.Errorf("getting summoner: %w", err)
	}

	err = pgxutil.QueryInsertRow(ctx, tx, "summoner_records", riotSummonerToRow(summoner))
	if err != nil {
		return fmt.Errorf("inserting summoner: %w", err)
	}

	leagues, err := db.riot.GetLeagueEntriesForSummoner(ctx, summoner.Id)
	if err != nil {
		return fmt.Errorf("getting league: %w", err)
	}

	row := map[string]any{"summoner_id": summoner.Id}
	for _, entry := range leagues {
		if entry.QueueType == riot.QueueTypeRankedSolo5x5 {
			for k, v := range riotLeagueEntryToRow(entry) {
				row[k] = v
			}
		}
	}

	err = pgxutil.QueryInsertRow(ctx, tx, "league_records", row)
	if err != nil {
		return fmt.Errorf("inserting league: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func riotSummonerToRow(m *riot.Summoner) map[string]any {
	return map[string]any{
		"summoner_id":     m.Id,
		"account_id":      m.AccountId,
		"puuid":           m.Puuid,
		"profile_icon_id": m.ProfileIconId,
		"summoner_level":  m.SummonerLevel,
		"revision_date":   riotUnixToDate(m.RevisionDate),
	}
}

func riotLeagueEntryToRow(m *riot.LeagueEntry) map[string]any {
	return map[string]any{
		"league_id":     m.LeagueId,
		"tier":          m.Tier,
		"division":      m.Rank,
		"league_points": m.LeaguePoints,
		"wins":          m.Wins,
		"losses":        m.Losses,
	}
}

type ProfileHeader struct {
	Name          string
	LastUpdated   string
	Rank          string
	WinLoss       string
	SummonerLevel int
}

func (db *DB) ProfileGetHeader(ctx context.Context, puuid string) (ProfileHeader, error) {
	rows, _ := db.pool.Query(ctx, `
	SELECT
		FORMAT('%s#%s', name, tagline) AS name,
		TO_CHAR(last_updated, 'YYYY MM-DD HH24:MI') AS last_updated,
		format_rank(tier, division, league_points) AS rank,
		format_win_loss(wins, losses) AS win_loss,
		summoner_level
	FROM profile_headers
	WHERE puuid = $1;
	`, puuid)

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[ProfileHeader])
}

type ProfileMatch struct {
	KillDeathAssist string
	CreepScore      string
	DamageDone      string
	GoldEarned      string
	VisionScore     string

	KillParticipation string
	CreepScorePer10   string
	DamagePercentage  string
	GoldPercentage    string

	MatchId      string
	GameDate     string
	GameDuration string
	PlayerWin    bool
	LpDelta      string

	ChampionIcon      string
	RunePrimaryIcon   string `db:"-"`
	RuneSecondaryIcon string `db:"-"`
	Spell1Icon        string `db:"-"`
	Spell2Icon        string `db:"-"`

	ItemIcons []*string
}

type ProfileMatchList []ProfileMatch

func (db *DB) ProfileGetMatchList(ctx context.Context, puuid string, page int, ensure bool) (ProfileMatchList, error) {
	start, count := 10*page, 10
	if ensure {
		err := ensureMatchList(ctx, db.pool, db.riot, puuid, start, count)
		if err != nil {
			return nil, fmt.Errorf("ensuring matchlist: %w", err)
		}
	}

	rows, _ := db.pool.Query(ctx, `
	SELECT
		FORMAT('%s / %s / %s', kills, deaths, assists) as kill_death_assist,
		creep_score,
		TO_CHAR(total_damage_dealt_to_champions, 'FM999,999') AS damage_done,
		TO_CHAR(gold_earned, 'FM999,999') AS gold_earned,
		vision_score,

		format_kill_participation() AS kill_participation,
		format_cs_per10(creep_score, game_duration) AS creep_score_per10,
		format_damage_relative() AS damage_percentage,
		'34%' AS gold_percentage,

		match_id,
		TO_CHAR(game_date, 'MM-DD HH24:MI') AS game_date,
		EXTRACT(MINUTE FROM game_duration) || ':' || EXTRACT(SECOND FROM game_duration) AS game_duration,
		win AS player_win,
		'??' AS lp_delta,

		get_champion_icon_url(champion_id) AS champion_icon,
		get_item_icon_urls(items) AS item_icons
	FROM profile_matches
	WHERE puuid = $1
	ORDER BY game_date DESC
	OFFSET $2 LIMIT $3;
	`, puuid, start, count)

	return pgx.CollectRows(rows, pgx.RowToStructByName[ProfileMatch])
}

type ProfileLiveGameParticipant struct {
	SummonerId        string
	SummonerName      string
	ChampionIcon      string
	RunePrimaryIcon   string
	RuneSecondaryIcon string
	Spell1Icon        string
	Spell2Icon        string
	Rank              string
	WinLoss           string
	WinLossRatio      string
}

type ProfileLiveGame struct {
	GameStartDate    string
	RedSide          []ProfileLiveGameParticipant
	BlueSide         []ProfileLiveGameParticipant
	RedSideBanIcons  []string
	BlueSideBanIcons []string
}

func (db *DB) ProfileGetLiveGame(ctx context.Context, name string) (ProfileLiveGame, error) {
	ids, err := db.GetAccount(ctx, name)
	if err != nil {
		return ProfileLiveGame{}, err
	}

	game, err := db.riot.GetCurrentGameInfoByPuuid(ctx, ids.Puuid)
	if err != nil {
		return ProfileLiveGame{}, err
	}

	var m ProfileLiveGame
	m.GameStartDate = riotUnixToDate(game.GameStartTime).String()

	for _, player := range game.Participants {
		name, err := riotGetName(ctx, db.riot, player.Puuid)
		if err != nil {
			return ProfileLiveGame{}, err
		}

		p := ProfileLiveGameParticipant{
			SummonerName: name,
			SummonerId:   player.SummonerId,
			// RunePrimaryIcon:   player.Perks.PerkStyle,
			RuneSecondaryIcon: ddragon.GetRuneTreeIconUrl(player.Perks.PerkSubStyle),
			Spell1Icon:        ddragon.GetSummonerSpellUrl(player.Spell1Id),
			Spell2Icon:        ddragon.GetSummonerSpellUrl(player.Spell2Id),
		}

		switch player.TeamId {
		case riot.TeamBlueSideID:
			m.BlueSide = append(m.BlueSide, p)
		case riot.TeamRedSideID:
			m.RedSide = append(m.RedSide, p)
		default:
			logging.FromContext(ctx).Sugar().DPanicf("invalid team id: %d", player.TeamId)
		}
	}

	m.BlueSideBanIcons = make([]string, 5)
	m.RedSideBanIcons = make([]string, 5)
	for _, ban := range game.BannedChampions {
		icon := ddragon.GetChampionIconUrl(1)
		switch ban.TeamId {
		case riot.TeamBlueSideID:
			m.BlueSideBanIcons[0] = icon
		case riot.TeamRedSideID:
			m.RedSideBanIcons[0] = icon
		default:
			logging.FromContext(ctx).Sugar().DPanicf("invalid team id: %d", ban.TeamId)
		}
	}

	return m, nil
}
