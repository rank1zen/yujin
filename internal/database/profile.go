package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/internal"
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

type ProfileHeader struct {
	Puuid         riot.PUUID `db:"puuid"`
	LastUpdated   time.Time  `db:"last_updated"`
	RiotID        string     `db:"name"`
	RiotTagLine   string     `db:"tagline"`
	SummonerLevel int        `db:"summoner_level"`

	RankWins      int    `db:"wins"`
	RankLosses    int    `db:"losses"`
	RankDivision  string `db:"division"`
	RankTier      int    `db:"tier"`
	RankLP        int    `db:"lp"`
	RankTimestamp int    `db:"rank_ts"`
}

func (db *DB) ProfileGetHeader(ctx context.Context, puuid string) (ProfileHeader, error) {
	rows, _ := db.pool.Query(ctx, `
	SELECT
		profile.puuid           AS puuid,
		profile.last_updated    AS last_updated,
		profile.name            AS name,
		profile.tagline         AS tagline,
		summoner.summoner_level AS summoner_level,
		league.wins             AS wins,
		league.losses           AS losses,
		league.division         AS division,
		league.tier             AS tier,
		league.league_points    AS lp,
		league.record_date 	AS rank_ts,
	FROM
		summoner_records_latest AS summoner
	JOIN
		league_records_latest AS league
		ON summoner.summoner_id = league.summoner_id
	JOIN
		profiles AS profile
		ON summoner.puuid = profile.puuid;
	`, puuid)
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[ProfileHeader])
}

type ProfileMatch struct {
	Puuid         riot.PUUID         `db:"puuid"`
	MatchID       riot.MatchID       `db:"match_id"`
	GameDate      time.Time          `db:"game_date"`
	GameDuration  time.Duration      `db:"game_duration"`
	GamePatch     string             `db:"game_patch"`
	Win           bool               `db:"win"`
	LpDelta       *int               `db:"lp_delta"`
	Champion      internal.Champion  `db:"champion"`
	ChampionLevel int                `db:"champion_level"`
	Summoners     internal.Summoners `db:"summoners"`
	Runes         internal.Runes     `db:"runes"`
	Items         internal.Items     `db:"items"`

	KdaKills         int     `db:"kda_kills"`
	KdaDeaths        int     `db:"kda_deaths"`
	KdaAssists       int     `db:"kda_assists"`
	KdaParticipation float32 `db:"kda_participation"`

	CsRaw   int     `db:"cs_raw"`
	CsPer10 float32 `db:"cs_per10"`

	DmgRaw            int     `db:"dmg_raw"`
	DmgPercentageTeam float32 `db:"dmg_percentage_team"`
	// DmgDeltaEnemy     int     `db:"dmg_delta_enemy"`

	GoldRaw            int     `db:"gold_raw"`
	GoldPercentageTeam float32 `db:"gold_percentage_team"`
	// GoldDeltaEnemy     int     `db:"gold_delta_enemy"`

	VisRaw int `db:"vis_raw"`
}

type ProfileMatchList struct {
	Puuid   riot.PUUID
	Page    int
	Count   int
	HasMore bool
	List    []ProfileMatch
}

func (db *DB) ProfileGetMatchList(ctx context.Context, puuid riot.PUUID, page int, ensure bool) (ProfileMatchList, error) {
	start, count := 10*page, 10
	if ensure {
		err := ensureMatchList(ctx, db.pool, db.riot, puuid, start, count)
		if err != nil {
			return ProfileMatchList{}, fmt.Errorf("ensuring matchlist: %w", err)
		}
	}

	rows, _ := db.pool.Query(ctx, `
	WITH
	team_total AS (
		SELECT
			team_id,
			match_id,
			sum(total_damage_dealt_to_champions) AS damage,
			sum(kills)                           AS kills,
			sum(gold_earned)                     AS gold
		FROM
			match_participants
		GROUP BY
			team_id, match_id
	)
	SELECT
		mp.puuid          AS puuid,
		mp.match_id       AS match_id,
		mi.game_date      AS game_date,
		mi.game_duration  AS game_duration,
		mi.game_patch     AS game_patch,
		mt.win 	          AS win,
		12 		  AS lp_delta,
		mp.champion_id    AS champion,
		mp.champion_level AS champion_level,
		mp.summoners      AS summoners,
		mp.runes          AS runes,
		mp.items          AS items,

		mp.kills                                                       AS kda_kills,
		mp.deaths                                                      AS kda_deaths,
		mp.assists                                                     AS kda_assists,
		round((mp.kills+mp.assists)/team_total.kills, 1)               AS kda_participation,
		mp.creep_score                                                 AS cs_raw,
		round(per_minute(mp.creep_score, mi.game_duration), 1)         AS cs_per10,
		mp.total_damage_dealt_to_champions                             AS dmg_raw,
		round(mp.total_damage_dealt_to_champions/team_total.damage, 2) AS dmg_percentage_team,
		mp.gold_earned                                                 AS gold_raw,
		round(mp.gold_earned/team_total.gold)                          AS gold_percentage_team,
		mp.vision_score                                                AS vis_raw
	FROM
		match_participants mp
	JOIN
		matches mi ON
			mi.id = mp.match_id
	JOIN
		match_teams mt ON 1=1
			AND mi.id = mt.match_id
			AND mp.team_id = mt.id
	JOIN
		team_total ON 1=1
			AND team_total.team_id = mp.team_id
			AND team_total.match_id = mp.match_id
	WHERE
		mp.puuid = $1
	ORDER BY
		game_date DESC
	OFFSET $2
	LIMIT $3;
	`, puuid, start, count)

	matches, err := pgx.CollectRows(rows, pgx.RowToStructByName[ProfileMatch])
	if err != nil {
		return ProfileMatchList{}, err
	}

	var hasMore bool
	db.pool.QueryRow(ctx, `
	SELECT EXISTS (
		SELECT
			1
		FROM
			match_participants mp
		WHERE
			mp.puuid = $1
		OFFSET $2
	);
	`, puuid, start).Scan(&hasMore)

	return ProfileMatchList{
		Puuid:   puuid,
		Page:    page,
		Count:   len(matches),
		HasMore: hasMore,
		List:    matches,
	}, nil
}

type ProfileLiveGameParticipant struct {
	SummonerID  riot.SummonerID `db:"-"`
	TeamID      riot.TeamID     `db:"-"`
	RiotID      string          `db:"name"`
	RiotTagLine string          `db:"tagline"`

	Champion  internal.Champion  `db:"-"`
	Summoners internal.Summoners `db:"-"`
	Runes     internal.Runes     `db:"-"`

	RankWins      int    `db:"wins"`
	RankLosses    int    `db:"losses"`
	RankDivision  string `db:"division"`
	RankTier      int    `db:"tier"`
	RankLP        int    `db:"lp"`
	RankTimestamp int    `db:"rank_ts"`

	BannedChampion *internal.Champion `db:"-"`
}

func scanLiveGameParticipants(m riot.SpectatorCurrentGameInfo, participants *[10]ProfileLiveGameParticipant) {
	for pID := range 10 {
		participants[pID].SummonerID = riot.SummonerID(m.Participants[pID].SummonerId)
		participants[pID].TeamID = riot.TeamID(m.Participants[pID].TeamId)
		participants[pID].Champion = internal.Champion(m.Participants[pID].ChampionId)
		participants[pID].Summoners = internal.Summoners{m.Participants[pID].Spell1Id, m.Participants[pID].Spell2Id}
		participants[pID].Runes = internal.Runes{m.Participants[pID].Perks.PerkIds[0]}
		participants[pID].BannedChampion = internal.GetBannedChampion(m.BannedChampions[pID].ChampionId)
	}
}

type ProfileLiveGame struct {
	StartDate time.Time

	Participants [10]ProfileLiveGameParticipant
}

func (db *DB) ProfileGetLiveGame(ctx context.Context, puuid string) (ProfileLiveGame, error) {
	m, err := db.riot.GetCurrentGameInfoByPuuid(ctx, puuid)
	if err != nil {
		return ProfileLiveGame{}, err
	}

	var participantPuuids []string
	for _, p := range m.Participants {
		participantPuuids = append(participantPuuids, p.Puuid)
	}

	rows, err := db.pool.Query(ctx, `
	WITH
	latest_rank AS (
		SELECT
			summoner_id,
			puuid
		FROM
			summoner_records_latest
		JOIN
			league_records_latest ON
				league_records_latest.summoner_id = summoner_records_latest.summoner_id
		WHERE
			puuid = $1
	)
	SELECT
		wins          AS wins,
		losses        AS losses,
		division      AS division,
		tier          AS tier,
		league_points AS lp
	FROM
		latest_rank
	JOIN
		profiles ON profiles.puuid = latest_rank.puuid
	WHERE
		puuid = ANY($1);
	`, participantPuuids)

	participants, err := pgx.CollectRows(rows, pgx.RowToStructByName[ProfileLiveGameParticipant])
	if err != nil {
		return ProfileLiveGame{}, err
	}

	p := [10]ProfileLiveGameParticipant(participants)

	scanLiveGameParticipants(m, &p)

	return ProfileLiveGame{
		StartDate:    riotUnixToDate(m.GameStartTime),
		Participants: p,
	}, nil
}

type ProfileChampionStat struct {
	Puuid       riot.PUUID        `db:"puuid"`
	Champion    internal.Champion `db:"champion"`
	GamesPlayed int               `db:"games_played"`
	LpDelta     int               `db:"-"`

	WlRaw        int    `db:"-"`
	WlPercentage string `db:"-"`

	KdaKills         int     `db:"kda_kills"`
	KdaDeaths        int     `db:"kda_deaths"`
	KdaAssists       int     `db:"kda_assists"`
	KdaParticipation float32 `db:"kda_participation"`

	CsRaw   int     `db:"cs_raw"`
	CsPer10 float32 `db:"cs_per10"`

	DmgRaw            int     `db:"dmg_raw"`
	DmgPercentageTeam float32 `db:"dmg_percentage_team"`
	// DmgDeltaEnemy     int     `db:"dmg_delta_enemy"`

	GoldRaw            int     `db:"gold_raw"`
	GoldPercentageTeam float32 `db:"gold_percentage_team"`
	// GoldDeltaEnemy     int     `db:"gold_delta_enemy"`

	VisRaw int `db:"vis_raw"`
}

type ProfileChampionStatList struct {
	Puuid  riot.PUUID
	Season internal.Season
	List   []ProfileChampionStat
}

func (db *DB) ProfileGetChampionStatList(ctx context.Context, puuid riot.PUUID, season internal.Season) (ProfileChampionStatList, error) {
	rows, _ := db.pool.Query(ctx, `
	WITH
	team_total AS (
		SELECT
			team_id,
			match_id,
			sum(total_damage_dealt_to_champions) AS damage,
			sum(kills)                           AS kills,
			sum(gold_earned)                     AS gold
		FROM
			match_participants
		GROUP BY
			team_id, match_id
	)
	participant AS (
		SELECT
			puuid,
			champion_id,
			count(*),
			sum(lp_delta),

			avg(mp.kills)                                                       AS kda_kills,
			avg(mp.deaths)                                                      AS kda_deaths,
			avg(mp.assists)                                                     AS kda_assists,
			avg(round(mp.kills/team_stats.kills, 2))                            AS kda_participation,
			avg(mp.creep_score)                                                 AS cs_raw,
			avg(mp.creep_score)                                                 AS cs_per10,
			avg(mp.total_damage_dealt_to_champions)                             AS dmg_raw,
			avg(round(mp.total_damage_dealt_to_champions/team_total.damage, 2)) AS dmg_percentage_team,
			avg(mp.gold_earned)                                                 AS gold_raw,
			avg()                                                               AS gold_percentage_team,
			avg(mp.vision_score)                                                AS vis_raw
		FROM
			match_participants mp
		JOIN
			team_stats
		WHERE
			puuid = $1
		GROUP BY
			champion_id
		ORDER BY
			champion_id
	)
	`, puuid, season)

	stats, err := pgx.CollectRows(rows, pgx.RowToStructByPos[ProfileChampionStat])
	if err != nil {
		return ProfileChampionStatList{}, err
	}

	return ProfileChampionStatList{
		Puuid:  puuid,
		Season: season,
		List:   stats,
	}, nil
}
