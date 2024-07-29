package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type SummonerMatchPostGame struct {
	Position      string `db:"player_position"`
	Items         []int  `db:"items_arr"`
	SummonerSpell []int  `db:"spells_arr"`
	Kills         int    `db:"kills"`
	Deaths        int    `db:"deaths"`
	Assists       int    `db:"assists"`
	CreepScore    int    `db:"creep_score"`
	VisionScore   int    `db:"vision_score"`
	ChampionLevel int    `db:"champion_level"`
	ChampionID    int    `db:"champion_id"`
	RunePrimary   int    `db:"rune_main_keystone"`
	RuneSecondary int    `db:"rune_secondary_path"`
}

func (m SummonerMatchPostGame) GetChampionIconUrl() string {
	return "https://ddragon.leagueoflegends.com/cdn/14.13.1/img/profileicon/871.png"
}

func (m SummonerMatchPostGame) GetSpellIconsUrls() []string {
	return []string{
		"https://ddragon.leagueoflegends.com/cdn/14.13.1/img/profileicon/871.png",
		"https://ddragon.leagueoflegends.com/cdn/14.13.1/img/profileicon/871.png",
	}
}

func (m SummonerMatchPostGame) GetItemIconUrls() []string {
	return []string{
		"https://ddragon.leagueoflegends.com/cdn/14.13.1/img/profileicon/871.png",
		"https://ddragon.leagueoflegends.com/cdn/14.13.1/img/profileicon/871.png",
		"https://ddragon.leagueoflegends.com/cdn/14.13.1/img/profileicon/871.png",
		"https://ddragon.leagueoflegends.com/cdn/14.13.1/img/profileicon/871.png",
		"https://ddragon.leagueoflegends.com/cdn/14.13.1/img/profileicon/871.png",
		"https://ddragon.leagueoflegends.com/cdn/14.13.1/img/profileicon/871.png",
	}
}

func (m SummonerMatchPostGame) GetRank() string {
	return "https://ddragon.leagueoflegends.com/cdn/14.13.1/img/profileicon/871.png"
}

func (m SummonerMatchPostGame) GetKda() string {
	return fmt.Sprintf("%d / %d / %d", m.Kills, m.Deaths, m.Assists)
}

func (m SummonerMatchPostGame) GetKills() string {
	return fmt.Sprintf("%d", m.Kills)
}

func (m SummonerMatchPostGame) GetDeaths() string {
	return fmt.Sprintf("%d", m.Deaths)
}

func (m SummonerMatchPostGame) GetAssists() string {
	return fmt.Sprintf("%d", m.Assists)
}

func (m SummonerMatchPostGame) GetKdaRatio() string {
	if m.Deaths == 0 {
		return "Perfect"
	}

	return fmt.Sprintf("%.2f", float32((m.Kills+m.Assists)/m.Deaths))
}

func (m SummonerMatchPostGame) GetDamage() string {
	return fmt.Sprintf("%dk", 1000)
}

func (m SummonerMatchPostGame) GetCreepScore() string {
	return fmt.Sprintf("%d", m.CreepScore)
}

func (m SummonerMatchPostGame) GetVisionScore() string {
	return fmt.Sprintf("%d", m.VisionScore)
}

type SummonerMatchPostGameList []*SummonerMatchPostGame

type ProfileSummary struct {
	LeagueTier     *string `db:"tier"`
	LeagueDivision *string `db:"division"`
	LeaguePoints   *int    `db:"league_points"`
	NumberWins     *int    `db:"number_wins"`
	NumberLosses   *int    `db:"number_losses"`
	ProfileIconId  int32   `db:"profile_icon_id"`
	SummonerLevel  int32   `db:"summoner_level"`
	Name           RiotName
}

func (m ProfileSummary) GetProfileIconUrl() string {
	return fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/14.13.1/img/profileicon/%d.png", m.ProfileIconId)
}

func (m ProfileSummary) GetSummonerLevel() string {
	return fmt.Sprintf("%d", m.SummonerLevel)
}

func (m ProfileSummary) GetWinLoss() string {
	if m.NumberWins == nil || m.NumberLosses == nil {
		return "0-0"
	}

	return fmt.Sprintf("%d-%d", *m.NumberWins, *m.NumberLosses)
}

func (m ProfileSummary) GetRank() string {
	if m.LeagueTier == nil {
		return "Unranked"
	}

	return fmt.Sprintf("%s %s", *m.LeagueTier, *m.LeagueDivision)
}

func (m ProfileSummary) GetLp() string {
	if m.LeaguePoints == nil {
		return ""
	}

	return fmt.Sprintf("%d LP", *m.LeaguePoints)
}

func (m ProfileSummary) GetName() string {
	return m.Name.String()
}

func (db *DB) GetProfileSummary(ctx context.Context, name RiotName) (*ProfileSummary, error) {
	ids, err := db.GetAccount(ctx, name)
	if err != nil {
		return nil, err
	}

	profile, err := db.getProfileSummary(ctx, ids.Puuid)
	if err != nil {
		return nil, fmt.Errorf("getProfileSummary: %w", err)
	}

	profile.Name = name
	return profile, nil
}

func (db *DB) getProfileSummary(ctx context.Context, puuid RiotPuuid) (*ProfileSummary, error) {
	sql := `
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
	`

	return querySelectRow(ctx, db.pool, sql, []any{puuid.String()}, pgx.RowToAddrOfStructByNameLax[ProfileSummary])
}

type ProfileMatch struct {
	MatchId      string        `db:"match_id"`
	GameDate     time.Time     `db:"game_date"`
	GameDuration time.Duration `db:"game_duration"`
	GamePatch    string        `db:"game_patch"`
	PlayerWin    bool          `db:"player_win"`
	SummonerMatchPostGame
}

func (m ProfileMatch) GetGamePatch() string {
	return m.GamePatch[:5]
}

func (m ProfileMatch) GetGameDate() string {
	duration := time.Since(m.GameDate)

	switch {
	case duration < time.Minute:
		return fmt.Sprintf("%d seconds ago", int(duration.Seconds()))
	case duration < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	case duration < time.Hour*24:
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	case duration < time.Hour*24*7:
		return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
	default:
		return m.GameDate.Format("2006-01-02")
	}
}

func (m ProfileMatch) GetGameDuration() string {
	return m.GameDuration.String()
}

type ProfileMatchList []*ProfileMatch

type ProfileMatchListPage struct {
	Page    int
	Count   int
	Name    RiotName
	Matches ProfileMatchList
}

func (p *ProfileMatchListPage) HasMore() bool {
	return true
}

func (db *DB) GetProfileMatchList(ctx context.Context, name RiotName, page int) (*ProfileMatchListPage, error) {
	ids, err := db.GetAccount(ctx, name)
	if err != nil {
		return nil, err
	}

	matches, err := db.getProfileMatchlist(ctx, ids.Puuid, page)
	if err != nil {
		return nil, fmt.Errorf("getProfileMatchlist: %w", err)
	}

	return &ProfileMatchListPage{
		Page:    page,
		Count:   len(matches),
		Name:    name,
		Matches: matches,
	}, nil
}

func (db *DB) getProfileMatchlist(ctx context.Context, puuid RiotPuuid, page int) (ProfileMatchList, error) {
	start, count := 10*page, 10

	rows, _ := db.pool.Query(ctx, `
	SELECT
		match_id,
		game_date,
		game_duration,
		game_patch,
		player_win,
		player_position,
		kills,
		deaths,
		assists,
		creep_score,
		champion_level,
		champion_id,
		vision_score,
		items_arr,
		spells_arr
	FROM
		match_summoner_postgame
	WHERE
		puuid = $1
	OFFSET $2 LIMIT $3;
	`, puuid.String(), start, count)

	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[ProfileMatch])
}

type ProfileMatchSummary struct {
	MatchId          string `db:"match_id"`
	PlayerSummaries  SummonerMatchPostGameList
	MatchItemEvents  []*MatchItemEvent
	MatchSkillEvents []*MatchSkillEvent
	MatchKillEvents  []*MatchKillEvent
}

type ProfileMatchSummaryList []*ProfileMatchSummary

// TODO
func (db *DB) GetProfileMatchSummary(ctx context.Context, puuid, matchID string) (*ProfileMatchSummary, error) {
	var m ProfileMatchSummary

	rows, _ := db.pool.Query(ctx, `
	SELECT
		player_position,
		kills,
		deaths,
		assists,
		creep_score,
		champion_level,
		champion_id,
		vision_score,
		items_arr,
		spells_arr,
	FROM
		match_participant_simple
	WHERE
		match_id = $1;
	`, matchID)

	players, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[SummonerMatchPostGame])
	if err != nil {
		return nil, err
	}

	if len(players) != 10 {
		return nil, fmt.Errorf("got players: %v", len(players))
	}

	m.PlayerSummaries = players

	rows, _ = db.pool.Query(ctx, `
	SELECT
		rune_slot,
		rune_id
	FROM
		match_rune_records
	WHERE
		puuid = $1 AND
		match_id = $2
	`, puuid, matchID)

	return &m, nil
}

func (db *DB) UpdateProfile(ctx context.Context, name RiotName) error {
	ids, err := db.GetAccount(ctx, name)
	if err != nil {
		return err
	}

	err = db.updateSummoner(ctx, ids.Puuid)
	if err != nil {
		return err
	}

	err = db.updateSummonerRankRecord(ctx, ids.SummonerId)
	if err != nil {
		return err
	}

	err = db.ensureMatchlist(ctx, ids.Puuid, 0, 20)
	if err != nil {
		return err
	}

	return nil
}
