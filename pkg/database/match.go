package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/pkg/riot"
)

type MatchInfo struct {
	RecordId     *string       `db:"record_id"`
	RecordDate   *time.Time    `db:"record_date"`
	MatchId      string        `db:"match_id"`
	GameDate     time.Time     `db:"game_date"`
	GameDuration time.Duration `db:"game_duration"`
	GamePatch    string        `db:"game_patch"`
}

type MatchTeam struct {
	TeamWin              bool `db:"team_win"`
	TeamSurrendered      bool `db:"team_surrendered"`
	TeamEarlySurrendered bool `db:"team_early_surrendered"`
}

type PlayerMetaInfo struct {
	PlayerId      int
	PlayerWin     bool   `db:"player_win"`
	Position      string `db:"player_position"`
	SummonerLevel int
}

type PlayerStats struct {
	Kills         int `db:"kills"`
	Deaths        int `db:"deaths"`
	Assists       int `db:"assists"`
	CreepScore    int `db:"creep_score"`
	GoldEarned    int `db:"gold_earned"`
	GoldSpent     int
	DoubleKills   int
	TripleKills   int
	QuadraKills   int
	PentaKills    int
	KillingSprees int

	ChampionExperience int
	ChampionLevel      int `db:"champion_level"`
	ChampionID         int `db:"champion_id"`
	ChampionTransform  int

	VisionScore             int
	WardsKilled             int
	WardsPlaced             int
	DetectorWardsPlaced     int
	SightWardsBoughtInGame  int
	VisionWardsBoughtInGame int
}

type PlayerDamageCharts struct {
	MagicDamageDealt               int
	MagicDamageDealtToChampions    int
	MagicDamageTaken               int
	PhysicalDamageDealt            int
	PhysicalDamageDealtToChampions int
	PhysicalDamageTaken            int
	TrueDamageDealt                int
	TrueDamageDealtToChampions     int
	TrueDamageTaken                int
	TotalDamageDealt               int
	TotalDamageDealtToChampions    int
	TotalDamageShieldedOnTeammates int
	TotalDamageTaken               int
	DamageDealtToBuildings         int
	DamageDealtToObjectives        int
	DamageDealtToTurrets           int
	DamageSelfMitigated            int
	TotalHeal                      int
	TotalHealsOnTeammates          int
}

type PlayerMiscStats struct {
	largestCriticalStrike         int
	largestKillingSpree           int
	largestMultiKill              int
	longestTimeSpentLiving        int
	neutralMinionsKilled          int
	objectivesStolen              int
	objectivesStolenAssists       int
	totalAllyJungleMinionsKilled  int
	totalEnemyJungleMinionsKilled int
	timeCCingOthers               int
	totalTimeCCDealt              int
	timePlayed                    int
	totalMinionsKilled            int
	totalTimeSpentDead            int
	totalUnitsHealed              int
	firstBloodAssist              bool
	firstBloodKill                bool
	firstTowerAssist              bool
	firstTowerKill                bool
	unrealKills                   int
	spell1Casts                   int
	spell2Casts                   int
	spell3Casts                   int
	spell4Casts                   int
	BaronKills                    int
	DragonKills                   int
	InhibitorKills                int
	NexusKills                    int
	TurretKills                   int
	TurretTakedowns               int
	NexusTakedowns                int
}

type PlayerPings struct {
	AllInPings         int
	AssistMePings      int
	CommandPings       int
	VisionClearedPings int
	DangerPings        int
	EnemyMissingPings  int
	EnemyVisionPings   int
	HoldPings          int
	GetBackPings       int
	NeedVisionPings    int
	OnMyWayPings       int
	PushPings          int
}

type PlayerItemRecord struct {
	RecordId string `db:"record_id"`
	MatchId  string `db:"match_id"`
	Puuid    string `db:"puuid"`
	ItemSlot int
	ItemID   int
}

type PlayerSpellRecord struct {
	RecordId  string `db:"record_id"`
	MatchId   string `db:"match_id"`
	Puuid     string `db:"puuid"`
	SpellSlot int
	SpellID   int
}

type MatchRuneRecord struct {
	RecordId string `db:"record_id"`
	MatchId  string `db:"match_id"`
	Puuid    string `db:"puuid"`

	Section   int
	Style     int
	Row       int
	Selection int
}

type MatchParticipantRecord struct {
	PlayerStats
	PlayerPings

	RecordId string `db:"record_id"`
	MatchId  string `db:"match_id"`
	Puuid    string `db:"puuid"`
	Level    int    `db:"level"`
}

type MatchObjectiveRecord struct {
	RecordId string `db:"record_id"`
	MatchId  string `db:"match_id"`
	TeamId   int32  `db:"team_id"`
	Name     string `db:"name"`
	First    bool   `db:"first"`
	Kills    int    `db:"kills"`
}

type MatchBanRecord struct {
	RecordId      string `db:"record_id"`
	MatchId       string `db:"match_id"`
	TeamId        int32  `db:"team_id"`
	BanChampionID int    `db:"champion_id"`
	BanTurn       int    `db:"turn"`
}

// UpdateMatchHistory fetches and inserts the most recent 10 matches
func (db *DB) UpdateMatchHistory(ctx context.Context, puuid string) ([]*MatchPlayer, error) {
	ids, err := db.riot.GetMatchHistory(ctx, puuid, 0, 10)
	if err != nil {
		return nil, err
	}

	matches, err := db.getMatchListPlayer(ctx, puuid, ids)
	if err != nil {
		return nil, err
	}

	return matches, nil
}

// GetMatchHistory fetches and inserts by page of 5 items
func (db *DB) GetMatchHistory(ctx context.Context, puuid string, page int) ([]*MatchPlayer, error) {
	start, count := 5*page, 5
	ids, err := db.riot.GetMatchHistory(ctx, puuid, start, count)
	if err != nil {
		return nil, err
	}

	matches, err := db.getMatchListPlayer(ctx, puuid, ids)
	if err != nil {
		return nil, err
	}

	return matches, nil
}

// getMatchListPlayer fetches each match corresponding to each matchIDs.
// Calls the Riot API and inserts matches as needed.
// Batch insert, nothing is inserted in case of failure.
// NOTE: This is currently not in transaction
func (db *DB) getMatchListPlayer(ctx context.Context, puuid string, matchIDs []string) ([]*MatchPlayer, error) {
	matches := make([]*MatchPlayer, len(matchIDs))

	for i, id := range matchIDs {
		match, err := db.getMatchPlayer(ctx, puuid, id)
		if err != nil {
			return nil, fmt.Errorf("failed %s: %w", id, err)
		}
		matches[i] = match
	}

	return matches, nil
}

func (db *DB) getMatchPlayer(ctx context.Context, puuid, matchID string) (*MatchPlayer, error) {
	batch := new(pgx.Batch)

	sql := `
	SELECT
		match_id, game_date, game_duration, game_patch,
		player_win, player_position, kills, deaths, assists, creep_score, champion_level, champion_id,
		items_arr,
		spells_arr
	FROM match_participant_simple
	WHERE match_id = $1 and puuid = $2;
	`

	row, _ := db.pool.Query(ctx, sql, matchID, puuid)
	match, err := pgx.CollectExactlyOneRow(row, pgx.RowToAddrOfStructByNameLax[MatchPlayer])

	var riotMatch *riot.Match

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		riotMatch, err = db.riot.GetMatch(ctx, matchID)
		if err != nil {
			return nil, fmt.Errorf("failed %s: %w", matchID, err)
		}
	case err != nil:
		return nil, fmt.Errorf("failed to check db: %w", err)
	default:
		return match, nil
	}

	batchMatch(batch, riotMatch)
	batch.Queue(sql, matchID, puuid)

	batchRes := db.pool.SendBatch(ctx, batch)
	defer batchRes.Close()

	for range batch.Len() - 1 {
		tag, err := batchRes.Exec()

		if err != nil {
			return nil, fmt.Errorf("batch insert: %v %w", tag, err)
		}
	}

	row, _ = batchRes.Query()
	match, err = pgx.CollectExactlyOneRow(row, pgx.RowToAddrOfStructByNameLax[MatchPlayer])
	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}

	return match, nil
}

func batchMatch(batch *pgx.Batch, m *riot.Match) {
	matchID := m.Metadata.MatchId

	var sql string
	var args []any

	sql, args = matchInfoQuery(m)
	batch.Queue(sql, args...)

	for _, p := range m.Info.Participants {
		sql, args = matchParticipantQuery(matchID, p)
		batch.Queue(sql, args...)
		sql, args = matchItemQuery(matchID, p)
		batch.Queue(sql, args...)
		sql, args = matchSummonerSpellQuery(matchID, p)
		batch.Queue(sql, args...)
	}

	for _, t := range m.Info.Teams {
		sql, args = matchTeamQuery(matchID, t)
		batch.Queue(sql, args...)
		sql, args = matchObjectiveQuery(matchID, t)
		batch.Queue(sql, args...)

		for _, ban := range t.Bans {
			sql, args = matchBanQuery(matchID, t.TeamId, ban)
			batch.Queue(sql, args...)
		}
	}
}

type MatchPlayer struct {
	MatchInfo
	PlayerMetaInfo
	PlayerStats

	// should have 6 items
	Items []int `db:"items_arr"`

	SummonerSpell []int `db:"spells_arr"`

	RunePrimaryID   int
	RuneSecondaryID int
}

func (m MatchPlayer) GetGameDate() string {
	return m.GameDate.String()
}

func (m MatchPlayer) GetGameDuration() string {
	return m.GameDuration.String()
}

func (m MatchPlayer) GetCreepScore() string {
	return "198"
}

func (m MatchPlayer) GetKDA() string {
	return fmt.Sprintf("%d/%d/%d", m.Kills, m.Deaths, m.Assists)
}

func (m MatchPlayer) GetKillDeathRatio() string {
	return fmt.Sprintf("%d", m.Kills+m.Assists/m.Deaths)
}

func (m MatchPlayer) GetChampionIconUrl() string {
	return "https://static.bigbrain.gg/assets/lol/riot_static/14.10.1/img/champion/Jhin.png"
}

func (m MatchPlayer) GetItemIconUrls() []string {
	return []string{
		"https://static.bigbrain.gg/assets/lol/riot_static/14.10.1/img/spell/SummonerFlash.png",
		"https://static.bigbrain.gg/assets/lol/riot_static/14.10.1/img/spell/SummonerFlash.png",
		"https://static.bigbrain.gg/assets/lol/riot_static/14.10.1/img/spell/SummonerFlash.png",
		"https://static.bigbrain.gg/assets/lol/riot_static/14.10.1/img/spell/SummonerFlash.png",
		"https://static.bigbrain.gg/assets/lol/riot_static/14.10.1/img/spell/SummonerFlash.png",
		"https://static.bigbrain.gg/assets/lol/riot_static/14.10.1/img/spell/SummonerFlash.png",
	}
}

func (m MatchPlayer) Valid() bool {
	return true
}

func (db *DB) getMatch(ctx context.Context, matchID string) (*MatchInfo, error) {
	rows, _ := db.pool.Query(ctx, `
	SELECT match_id, m.game_date, m.game_duration, m.game_patch
	FROM MatchInfoRecords AS m
	WHERE m.match_id = $1;
	`, matchID)

	match, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[MatchInfo])
	if err != nil {
		return nil, err
	}

	return match, nil
}

// This is probably dead
// Return new match ids (matches not currently in db) from Riot
//
// TODO: This could be made infinitely better
func (db *DB) fetchNewMatches(ctx context.Context, puuid string, start, count int) ([]string, error) {
	matchIDs, err := db.riot.GetMatchHistory(ctx, puuid, start, count)
	if err != nil {
		return nil, fmt.Errorf("failed to get matchlist ids from riot: %w", err)
	}

	rows, err := db.pool.Query(ctx, `
		SELECT match_id FROM match_info_records WHERE match_id = ANY ($1);
	`, matchIDs)
	if err != nil {
		return nil, fmt.Errorf("something went wrong with sql query: %w", err)
	}

	foundIDs, err := pgx.CollectRows(rows, pgx.RowToStructByName[string])
	if err != nil {
		return nil, fmt.Errorf("collect rows: %w", err)
	}

	newIDs := make([]string, 0)
	for _, id := range matchIDs {
		found := false
		for _, i := range foundIDs {
			if i == id {
				found = true
			}
		}
		if !found {
			newIDs = append(newIDs, id)
		}
	}

	return newIDs, nil
}
