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

// Returns true if profile with puuid exists. This function will fetch from riot.
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

type RiotID struct {
	Name string
	Tag  string
}

type ProfileHeader struct {
	Puuid         riot.PUUID     `db:"puuid"`
	Name          RiotID         `db:"name"`
	LastUpdated   time.Time      `db:"last_updated"`
	Rank          *RankTimestamp `db:"rank"`
	SummonerLevel int            `db:"summoner_level"`
}

func (db *DB) ProfileGetHeader(ctx context.Context, puuid string) (ProfileHeader, error) {
	rows, _ := db.pool.Query(ctx, `
	SELECT
		puuid,
		comp(name, tagline) AS name,
		last_updated,
		comp(tier, division) AS rank,
		summoner_level
	FROM profile_headers
	WHERE puuid = $1;
	`, puuid)
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[ProfileHeader])
}

type ParticipantItems [7]*int

type ParticipantRunes [9]int

type ParticipantSummoners [2]int

type ParticipantChampion int

type ProfileMatch struct {
	Puuid   riot.PUUID   `db:"puuid"`
	MatchID riot.MatchID `db:"match_id"`

	Kills        int           `db:"kills"`
	Deaths       int           `db:"deaths"`
	Assists      int           `db:"assists"`
	CreepScore   int           `db:"creep_score"`
	DamageDealt  int           `db:"damage_dealt"`
	GoldEarned   int           `db:"gold_earned"`
	VisionScore  int           `db:"vision_score"`
	GameDate     time.Time     `db:"game_date"`
	GameDuration time.Duration `db:"game_duration"`
	PlayerWin    bool          `db:"player_win"`
	LpGain       int           `db:"-"`

	DamagePercentage  float32 `db:"-"`
	GoldPercentage    float32 `db:"-"`
	KillParticipation float32 `db:"-"`

	Champion  ParticipantChampion  `db:"-"`
	Summoners ParticipantSummoners `db:"-"`
	Runes     ParticipantRunes     `db:"-"`
	Items     ParticipantItems     `db:"-"`
}

type ProfileMatchList struct {
	Page    int
	Count   int
	HasMore bool
	Puuid   riot.PUUID
	M       []ProfileMatch
}

func (db *DB) ProfileGetMatchList(ctx context.Context, puuid string, page int, ensure bool) (ProfileMatchList, error) {
	start, count := 10*page, 10
	if ensure {
		err := ensureMatchList(ctx, db.pool, db.riot, puuid, start, count)
		if err != nil {
			return ProfileMatchList{}, fmt.Errorf("ensuring matchlist: %w", err)
		}
	}

	var m ProfileMatchList

	rows, _ := db.pool.Query(ctx, `
	SELECT
		match_id                        AS match_id,
		kills                           AS kills,
		deaths                          AS deaths,
		assists                         AS assists,
		creep_score                     AS creep_score,
		total_damage_dealt_to_champions AS damage_dealt,
		gold_earned                     AS gold_earned,
		vision_score                    AS vision_score,
		game_date                       AS game_date,
		game_duration                   AS game_duration,
		win                             AS win,

		get_champion_icon_url(champion_id)
		get_rune_icon_url(rune_primary_keystone)
		get_rune_tree_icon_url(rune_secondary_path)
		get_summoners_icon_urls(summoners)
		get_item_icon_urls(items)
	FROM profile_matches
	WHERE puuid = $1
	ORDER BY game_date DESC
	OFFSET $2 LIMIT $3;
	`, puuid, start, count)

	a, err := pgx.CollectRows(rows, pgx.RowToStructByName[ProfileMatch])

	m.M = a

	return m, err
}

type Rank struct {
	Division string
	Tier     int
	LP       int
}

type RankTimestamp struct {
	Rank
	Wins      int
	Losses    int
	Timestamp time.Time
}

type ProfileLiveGameParticipant struct {
	SummonerID riot.SummonerID      `db:"-"`
	TeamID     riot.TeamID          `db:"-"`
	Name       string               `db:"-"`
	Rank       *RankTimestamp       `db:"-"`
	Champion   ParticipantChampion  `db:"-"`
	Summoners  ParticipantSummoners `db:"-"`
	Runes      ParticipantRunes     `db:"-"`
}

type ProfileLiveGameTeam struct {
	TeamID             riot.TeamID
	AverageRank        *Rank
	BannedChampionIcon [5]*string
	Participants       [5]ProfileLiveGameParticipant
}

type ProfileLiveGame struct {
	StartDate time.Time
	RedSide   ProfileLiveGameTeam
	BlueSide  ProfileLiveGameTeam
}

func (db *DB) ProfileGetLiveGame(ctx context.Context, puuid string) (ProfileLiveGame, error) {
	game, err := db.riot.GetCurrentGameInfoByPuuid(ctx, puuid)
	if err != nil {
		return ProfileLiveGame{}, err
	}

	var m ProfileLiveGame
	m.StartDate = riotUnixToDate(game.GameStartTime)
	for _, ban := range game.BannedChampions {
		if ban.TeamId == 1 {
		} else {
		}
	}

	for i, participant := range game.Participants {
		var p *ProfileLiveGameParticipant
		if participant.TeamId == 1 {
			p = &m.BlueSide.Participants[i]
		}

		p.Champion = ParticipantChampion(participant.ChampionId)
	}

	rows, err := db.pool.Query(ctx, `
	SELECT
		row('a', 'a')
	WHERE puuid = $1;
	`, puuid)

	for rows.Next() {
		rows.Scan()
	}

	return o, nil
}

type ProfileChampionStat struct {
	ChampionIcon      string
	LpDelta           string `db:"-"`
	GamesPlayed       string
	WinLoss           string `db:"-"`
	WinRate           string `db:"-"`
	KillDeathAssist   string `db:"-"`
	KillParticipation string `db:"-"`
	CreepScore        string
	CreepScorePer10   string `db:"-"`
	DamageDone        string
	DamagePercentage  string `db:"-"`
	GoldEarned        string
	GoldPercentage    string `db:"-"`
	VisionScore       string
}

type ProfileChampionStatList struct {
	Season string
	Stats  []ProfileChampionStat
}

func (db *DB) ProfileGetChampionStatList(ctx context.Context, puuid string, season string) (ProfileChampionStatList, error) {
	rows, _ := db.pool.Query(ctx, `
	WITH champ_stats AS (
		SELECT
			get_champion_icon_url(champion_id),
			count(*),
			avg(creep_score),
			avg(vision_score),
			avg(total_damage_dealt_to_champions),
			avg(gold_earned)
		FROM profile_matches
		WHERE puuid = $1
		GROUP BY champion_id
	)

	SELECT $2, array(SELECT row(champ_stats.*) FROM champ_stats);
	`, puuid, season)

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[ProfileChampionStatList])
}
