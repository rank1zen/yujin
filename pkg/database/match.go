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
func getPlayerHistory(ctx context.Context, db pgxDB, puuid string, start int, count int) ([]MatchCard, error) {
	rows, _ := db.Query(ctx, `
	SELECT
		m.match_id, m.game_date, m.game_duration, m.game_patch, p.player_win,
		p.player_position, p.kills, p.deaths, p.assists, p.creep_score,
		p.champion_level, p.champion_id
	FROM MatchInfoRecords AS m
	INNER JOIN MatchParticipantRecords AS p ON m.match_id = p.match_id
	WHERE p.puuid = $1
	ORDER BY m.game_date
	OFFSET $2 LIMIT $3;
	`, puuid, start, count)

	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[MatchCard])
}

// this is supposed to find new matches from riot not in the sysmtem
func checkNewMatchlist(ctx context.Context, db pgxDB, riot RiotClient, puuid string) ([]string, error) {
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

func updateMatchlist(ctx context.Context, db pgxDB, riot RiotClient, puuid string) error {
	logger := logging.FromContext(ctx).Sugar()

	newIDs, err := checkNewMatchlist(ctx, db, riot, puuid)
	if err != nil {
		return err
	}

	logger.Debugf("fetching match ids: %v", newIDs)
	// TODO: probably add some cancel mechanism here
	err = insertMatches(ctx, db, riot, newIDs)
	if err != nil {
		return fmt.Errorf("failed to fetch each match: %w", err)
	}

	return nil
}

func insertMatches(ctx context.Context, db pgxDB, riot RiotClient, matchIDs []string) error {
	return pgx.BeginFunc(ctx, db, func(tx pgx.Tx) error {
		for _, id := range matchIDs {
			m, err := riot.GetMatch(id)
			if err != nil {
				return fmt.Errorf("failed: %s from riot: %w", id, err)
			}

			_, err = tx.Exec(ctx, `
			INSERT INTO MatchInfoRecords
			(match_id, game_date, game_duration, game_patch)
			VALUES ($1, $2, $3, $4);
			`, m.Metadata.MatchID, m.Info.GameStartTimestamp, m.Info.GameDuration,
				m.Info.GameVersion)
			if err != nil {
				return err
			}

			for _, p := range m.Info.Participants {
				_, err = tx.Exec(ctx, `
				INSERT INTO MatchParticipantRecords
				(match_id, puuid, player_win, player_position, kills,
				deaths, assists, creep_score, gold_earned, champion_level,
				champion_id)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
				`, m.Metadata.MatchID, p.PUUID, p.Win, p.Role, p.Kills,
					p.Deaths, p.Assists, p.PUUID, p.GoldEarned, p.ChampLevel,
					p.ChampionID)
				if err != nil {
					return err
				}
			}

			for _, t := range m.Info.Teams {
				_, err = tx.Exec(ctx, `
				INSERT INTO MatchTeamRecords
				(match_id, team_id, team_win, team_surrendered, team_early_surrendered)
				VALUES ($1, $2, $3, $4, $5)
				`, m.Metadata.MatchID, t.TeamID, t.Win, false, false) // FIXME
				if err != nil {
					return err
				}
			}

		}

		return nil
	})
}
