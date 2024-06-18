package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/pkg/logging"
)

type MatchInfo struct {
	MatchId      string        `db:"match_id"`
	GameDate     time.Time     `db:"game_date"`
	GameDuration time.Duration `db:"game_duration"`
	GamePatch    string        `db:"game_patch"`
}

type MatchInfoRecord struct {
	RecordId   string    `db:"record_id"`
	RecordDate time.Time `db:"record_date"`

	MatchInfo
}

type MatchTeam struct {
	TeamWin              bool `db:"team_win"`
	TeamSurrendered      bool `db:"team_surrendered"`
	TeamEarlySurrendered bool `db:"team_early_surrendered"`
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

type PlayerItem struct {
	RecordId string `db:"record_id"`
	MatchId  string `db:"match_id"`
	Puuid    string `db:"puuid"`
	ItemSlot int
	ItemID   int
}

type PlayerSpell struct {
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

type MatchCard struct {
	MatchInfo
	PlayerMetaInfo
	PlayerStats

	Items []PlayerItem

	SummonerSpell1 int
	SummonerSpell2 int

	RunePrimary   int
	RuneSecondary int
}

// gets MatchInfo, MatchParticipant (just the player), and the all players items, runes, summoner spells
func getPlayerMatchHistory(ctx context.Context, db pgxDB, puuid string, start int, count int) ([]MatchCard, error) {
	rows, _ := db.Query(ctx, `
	SELECT
		m.match_id, m.game_date, m.game_duration, m.game_patch, p.player_win,
		p.player_position, p.kills, p.deaths, p.assists, p.creep_score,
		p.champion_level, p.champion_id
	FROM MatchInfoRecords AS m
	INNER JOIN MatchParticipantRecords AS p ON m.match_id = p.match_id
	WHERE p.puuid = $1
	ORDER BY m.game_date DESC
	OFFSET $2 LIMIT $3;
	`, puuid, start, count)

	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[MatchCard])
}

// Return new match ids from Riot
func fetchMatchlist(ctx context.Context, db pgxDB, riot RiotClient, puuid string) ([]string, error) {
	matchIDs, err := riot.GetMatchlist(puuid, 0, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get matchlist ids from riot: %w", err)
	}

	rows, err := db.Query(ctx, `
		SELECT match_id FROM MatchInfoRecords WHERE match_id = ANY ($1);
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

// Iterates over each id, fetches and inserts all match records, returns the inserted IDs
// Batch insert, nothing is inserted in case of failure.
// If the database has the match then we ignore
func insertMatches(ctx context.Context, db pgxDB, riot RiotClient, matchIDs []string) ([]string, error) {
	batch := &pgx.Batch{}
	insertedIDs := make([]string, 0)

	for _, id := range matchIDs {
		var f bool
		err := db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM MatchInfoRecords WHERE match_id = $1)`, id).Scan(&f)
		if err != nil {
			return nil, fmt.Errorf("failed to check db: %w", err)
		}

		if f {
			continue
		}

		m, err := riot.GetMatch(id)
		if err != nil {
			return nil, fmt.Errorf("failed: %s from riot: %w", id, err)
		}

		gameDate := time.Unix(m.Info.GameStartTimestamp, 0)
		gameDuration := time.Duration(m.Info.GameDuration) * time.Millisecond

		batch.Queue(`
		INSERT INTO MatchInfoRecords
		(match_id, game_date, game_duration, game_patch)
		VALUES ($1, $2, $3, $4);
		`, m.Metadata.MatchID, gameDate, gameDuration, m.Info.GameVersion)

		for _, p := range m.Info.Participants {
			batch.Queue(`
			INSERT INTO MatchParticipantRecords
			(match_id, puuid, player_win, player_position, kills,
			deaths, assists, creep_score, gold_earned, champion_level,
			champion_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
			`, m.Metadata.MatchID, p.PUUID, p.Win, p.Role, p.Kills,
				p.Deaths, p.Assists, p.TotalMinionsKilled, p.GoldEarned, p.ChampLevel,
				p.ChampionID)
		}

		for _, t := range m.Info.Teams {
			batch.Queue(`
			INSERT INTO MatchTeamRecords
			(match_id, team_id, team_win, team_surrendered, team_early_surrendered)
			VALUES ($1, $2, $3, $4, $5)
			`, m.Metadata.MatchID, t.TeamID, t.Win, false, false) // FIXME
		}

		insertedIDs = append(insertedIDs, id)
	}

	br := db.SendBatch(ctx, batch)

	_, err := br.Exec()
	if err != nil {
		return nil, fmt.Errorf("batch failed: %w", err)
	}

	return insertedIDs, nil
}

// Calls fetchMatchlist to fetches match history then call insertMatches on new ids
func updateMatchHistory(ctx context.Context, db pgxDB, riot RiotClient, puuid string) ([]string, error) {
	logger := logging.FromContext(ctx).Sugar()

	newIDs, err := fetchMatchlist(ctx, db, riot, puuid)
	if err != nil {
		return nil, err
	}

	logger.Debugf("fetching match ids: %v", newIDs)
	// TODO: probably add some cancel mechanism here
	insertedIDs, err := insertMatches(ctx, db, riot, newIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch each match: %w", err)
	}

	return insertedIDs, nil
}
