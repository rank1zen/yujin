package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
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
	Puuid         string
	Name          string
	LastUpdated   string
	Rank          string
	WinLoss       string
	SummonerLevel int
}

func (db *DB) ProfileGetHeader(ctx context.Context, puuid string) (ProfileHeader, error) {
	rows, _ := db.pool.Query(ctx, `
	SELECT
		puuid,
		format('%s#%s', name, tagline) as name,
		to_char(last_updated, 'YYYY MM-DD HH24:MI') AS last_updated,
		format_rank(tier, division, league_points) AS rank,
		format_win_loss(wins, losses) AS win_loss,
		summoner_level
	FROM profile_headers
	WHERE puuid = $1;
	`, puuid)
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[ProfileHeader])
}

type ProfileMatch struct {
	MatchID           riot.MatchID
	KillDeathAssist   string
	KillParticipation string
	CreepScore        string
	CreepScorePer10   string
	DamageDone        string
	DamagePercentage  string
	GoldEarned        string
	GoldPercentage    string
	VisionScore       string
	GameDate          string
	GameDuration      string
	PlayerWin         bool
	LpDelta           string

	ChampionIcon      string
	RunePrimaryIcon   string
	RuneSecondaryIcon string
	SummonersIcons    [2]string
	ItemIcons         [7]*string
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
		match_id,
		format('%s / %s / %s', kills, deaths, assists) AS kill_death_assist,
		format_kill_participation() AS kill_participation,
		creep_score,
		format_cs_per10(creep_score, game_duration) AS creep_score_per10,
		to_char(total_damage_dealt_to_champions, 'FM999,999') AS damage_done,
		format_damage_relative() AS damage_percentage,
		to_char(gold_earned, 'FM999,999') AS gold_earned,
		'34%' AS gold_percentage,
		vision_score,
		to_char(game_date, 'MM-DD HH24:MI') AS game_date,
		extract(minute from game_duration) || ':' || extract(second from game_duration) AS game_duration,
		win AS player_win,
		'??' AS lp_delta,

		get_champion_icon_url(champion_id) AS champion_icon,
		get_rune_icon_url(rune_primary_keystone) AS rune_primary_icon,
		get_rune_tree_icon_url(rune_secondary_path) AS rune_secondary_icon,
		get_summoners_icon_urls(summoners) AS summoners_icons,
		get_item_icon_urls(items) AS item_icons
	FROM profile_matches
	WHERE puuid = $1
	ORDER BY game_date DESC
	OFFSET $2 LIMIT $3;
	`, puuid, start, count)

	return pgx.CollectRows(rows, pgx.RowToStructByName[ProfileMatch])
}

type ProfileLiveGameParticipant struct {
	SummonerID         string
	TeamID             riot.TeamID
	Name               string
	Rank               string
	WinLoss            string
	ChampionBannedIcon *string
	ChampionIcon       string
	SummonersIcons     [2]string
	RunePrimaryIcon    string
	RuneSecondaryIcon  string
}

type ProfileLiveGameTeam struct {
	Participants [5]ProfileLiveGameParticipant
}

type ProfileLiveGame struct {
	GameStartDate string
	BlueTeam      ProfileLiveGameTeam
	RedTeam       ProfileLiveGameTeam
}

func (db *DB) ProfileGetLiveGame(ctx context.Context, puuid string) (ProfileLiveGame, error) {
	m, err := db.riot.GetCurrentGameInfoByPuuid(ctx, puuid)
	if err != nil {
		return ProfileLiveGame{}, err
	}

	// HACK: we want to do this through sql instead of go
	var livegame ProfileLiveGame
	livegame.GameStartDate = riotUnixToDate(m.GameStartTime).String()

	for i, player := range m.Participants {
		p, err := riotSpectatorParticipant(ctx, db.pool, db.riot, player)
		if err != nil {
			return ProfileLiveGame{}, err
		}
		// ULTRA HACK
		if p.TeamID == 100 {
			livegame.RedTeam.Participants[i%5] = p
		} else {
			livegame.BlueTeam.Participants[i%5] = p
		}
	}

	return livegame, nil
}

func riotSpectatorParticipant(ctx context.Context, db pgxutil.Conn, r *riot.Client, m riot.SpectatorCurrentGameParticipant) (ProfileLiveGameParticipant, error) {
	name, err := riotGetName(ctx, r, m.Puuid)
	if err != nil {
		return ProfileLiveGameParticipant{}, err
	}

	summonersIcons, err := dbGetSummonersIconUrls(ctx, db, [2]int{m.Spell1Id, m.Spell2Id})
	if err != nil {
		return ProfileLiveGameParticipant{}, err
	}

	championIcon, err := dbGetChampionIconUrl(ctx, db, m.ChampionId)
	if err != nil {
		return ProfileLiveGameParticipant{}, err
	}

	runePrimaryIcon, err := dbGetRuneIconUrl(ctx, db, m.Perks.PerkIds[riot.PerkKeystone])
	if err != nil {
		return ProfileLiveGameParticipant{}, err
	}

	runeSecondaryIcon, err := dbGetRuneTreeIconUrl(ctx, db, m.Perks.PerkStyle)
	if err != nil {
		return ProfileLiveGameParticipant{}, err
	}

	return ProfileLiveGameParticipant{
		SummonerID:         m.SummonerId,
		TeamID:             riot.TeamID(m.TeamId),
		Name:               name,
		Rank:               "",
		WinLoss:            "",
		ChampionBannedIcon: nil,
		ChampionIcon:       championIcon,
		RunePrimaryIcon:    runePrimaryIcon,
		RuneSecondaryIcon:  runeSecondaryIcon,
		SummonersIcons:     summonersIcons,
	}, nil
}
