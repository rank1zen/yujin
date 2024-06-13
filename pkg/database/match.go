package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type MatchInfoRecord struct {
	RecordId   string        `db:"record_id"`
	RecordDate time.Time     `db:"record_date"`
	MatchId    string        `db:"match_id"`
	Patch      string        `db:"game_patch"`
	Duration   time.Duration `db:"game_duration"`
	Date       time.Time     `db:"game_date"`
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
	ItemID int
}

type PlayerSpell struct {
	Spell int
}

type MatchRune struct {
	Curr          int
	RuneSelection int
	RuneRow       int
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
	PlayerMetaInfo
	PlayerStats

	Items []PlayerItem

	SummonerSpell1 int
	SummonerSpell2 int

	RunePrimary   int
	RuneSecondary int
}

func getMatchHistory(ctx context.Context, db pgxDB, puuid string) ([]MatchCard, error) {
	rows, _ := db.Query(ctx, `
		SELECT
			t1.duration, t1.patch, t1.start_ts
		FROM MatchInfoRecords AS t1
		JOIN (
			SELECT 
			FROM MatchParticipantRecords
			WHERE puuid = $1
		) AS t2
			ON t2.match_id = t1.match_id
		ORDER BY record_date DESC;
	`, puuid)

	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[MatchCard])
}

func fetchAndInsertMatches(ctx context.Context, db pgxDB, riot RiotClient, puuid string) error {
	ids, err := riot.GetMatchlist(puuid)
	if err != nil {
		return err
	}

	rows, _ := db.Query(ctx, `
		SELECT match_id FROM MatchRecords WHERE match_id IN ANY($1)
	`, ids)

	newIDs, err := pgx.CollectRows(rows, pgx.RowToStructByName[string])
	if err != nil {
		return err
	}

	err = pgx.BeginFunc(ctx, db, func(tx pgx.Tx) error {
		for _, id := range newIDs {
			m, err := riot.GetMatch(id)
			if err != nil {
				return err
			}

			_, err = tx.Exec(ctx, `
				INSERT INTO MatchInfoRecords
				(match_id, game_date, game_duration, game_patch)
				VALUES ($1, $2, $3, $4)
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
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
				`, m.Metadata.MatchID, p.PUUID, p.Win, p.Role, p.Kills,
				p.Deaths, p.Assists, p.PUUID, p.GoldEarned, p.ChampLevel,
				p.ChampionID)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
