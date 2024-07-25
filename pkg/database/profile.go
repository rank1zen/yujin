package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type SummonerMatchPostGame struct {
	Kills         int    `db:"kills"`
	Deaths        int    `db:"deaths"`
	Assists       int    `db:"assists"`
	CreepScore    int    `db:"creep_score"`
	VisionScore   int    `db:"vision_score"`
	ChampionLevel int    `db:"champion_level"`
	ChampionID    int    `db:"champion_id"`
	Position      string `db:"player_position"`

	Items         []int `db:"items_arr"`  // should have 6 items
	SummonerSpell []int `db:"spells_arr"` // should have 2
	RunePrimary   int   `db:"main_keystone"`
	RuneSecondary int   `db:"secondary_path"`
}

func (m SummonerMatchPostGame) GetChampionIconUrl() string {
	panic("jerry")
}

func (m SummonerMatchPostGame) GetSpellIconsUrls() []string {
	panic("jerry")
}

func (m SummonerMatchPostGame) GetItemIconUrls() []string {
	panic("jerry")
}

func (m SummonerMatchPostGame) GetRank() string {
	panic("jerry")
}

func (m SummonerMatchPostGame) GetKda() string {
	return fmt.Sprintf("%d / %d / %d", m.Kills, m.Deaths, m.Assists)
}

func (m SummonerMatchPostGame) GetKdaRatio() string {
	if m.Deaths == 0 {
		return "Perfect KDA"
	}

	return fmt.Sprintf("%.2f KDA", float32((m.Kills+m.Assists)/m.Deaths))
}

func (m SummonerMatchPostGame) GetDamage() string {
	return fmt.Sprintf("%dk", 1000)
}

func (m SummonerMatchPostGame) GetCS() string {
	return fmt.Sprintf("%d", m.CreepScore)
}

type SummonerMatchPostGameList []*SummonerMatchPostGame

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

func (m ProfileMatch) GetCreepScore() string {
	return fmt.Sprintf("%d CS", m.CreepScore)
}

func (m ProfileMatch) GetKDA() string {
	return fmt.Sprintf("%d / %d / %d", m.Kills, m.Deaths, m.Assists)
}

func (m ProfileMatch) GetKillDeathRatio() string {
	if m.Deaths == 0 {
		return "Perfect KDA"
	}

	return fmt.Sprintf("%.2f KDA", float32((m.Kills+m.Assists)/m.Deaths))
}

func (m ProfileMatch) GetVisionScore() string {
	return fmt.Sprintf("%d", m.VisionScore)
}

func (m ProfileMatch) GetChampionIconUrl() string {
	return "https://static.bigbrain.gg/assets/lol/riot_static/14.10.1/img/champion/Jhin.png"
}

func (m ProfileMatch) GetSpellsIconUrls() []string {
	return []string{
		"https://static.bigbrain.gg/assets/lol/riot_static/14.10.1/img/spell/SummonerFlash.png",
		"https://static.bigbrain.gg/assets/lol/riot_static/14.10.1/img/spell/SummonerFlash.png",
	}
}

func (m ProfileMatch) GetItemIconUrls() []string {
	urls := make([]string, 0)
	for id := range m.Items {
		item := fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/14.13.1/img/item/%d.png", id)
		urls = append(urls, item)
	}

	return urls
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
		return nil, err
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
		spells_arr,
		runes_arr
	FROM match_participant_simple
	WHERE puuid = $1
	OFFSET $2 LIMIT $3
	`, puuid, start, count)

	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[ProfileMatch])
}

type ProfileMatchSummary struct {
	MatchId          string `db:"match_id"`
	PlayerSummaries  SummonerMatchPostGameList
	MatchItemEvents  []*MatchItemEvent
	MatchSkillEvents []*MatchSkillEvent
	MatchKillEvents  []*MatchKillEvent
	Runes            *MatchRuneFull
}

type ProfileMatchSummaryList []*ProfileMatchSummary

type ProfileSummary struct {
	ProfileIconId  int32   `db:"profile_icon_id"`
	SummonerLevel  int32   `db:"summoner_level"`
	LeagueTier     *string `db:"tier"`
	LeagueDivision *string `db:"division"`
	LeaguePoints   *string `db:"league_points"`
	NumberWins     *int    `db:"number_wins"`
	NumberLosses   *int    `db:"number_losses"`
}

func (m ProfileSummary) GetName() string {
	return "Doublelift"
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

func (m ProfileSummary) GetRankLP() string {
	if m.LeaguePoints == nil {
		return "0 LP"
	}

	return fmt.Sprintf("%d LP", *m.LeaguePoints)
}

func (db *DB) getProfileSummary(ctx context.Context, puuid string) (*ProfileSummary, error) {
	return SelectRow(ctx, db.pool, `
	SELECT
		profile_icon_id,
		summoner_level,
		tier,
		division,
		league_points,
		number_wins,
		number_losses
	FROM summoner_profile
	WHERE puuid = $1;
	`, []any{
		puuid,
	}, pgx.RowToAddrOfStructByName[ProfileSummary])
}

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

	var runeSlot RuneSlot
	var runeID int
	_, err = pgx.ForEachRow(rows, []any{&runeSlot, &runeID}, func() error {
		return identifyRuneFull(m.Runes, runeSlot, runeID)
	})
	if err != nil {
		return nil, err
	}

	return &m, nil
}
