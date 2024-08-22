package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/pkg/logging"
	"github.com/rank1zen/yujin/pkg/riot"
)

type MatchSummonerPostGame struct {
	Item0Id int `db:"item0_id"`
	Item1Id int `db:"item1_id"`
	Item2Id int `db:"item2_id"`
	Item3Id int `db:"item3_id"`
	Item4Id int `db:"item4_id"`
	Item5Id int `db:"item5_id"`
	Item6Id int `db:"item6_id"`

	Spells0Id int `db:"spell0_id"`
	Spells1Id int `db:"spell1_id"`

	RunePrimaryPath     int `db:"rune_primary_path"`
	RunePrimaryKeystone int `db:"rune_primary_keystone"`
	RunePrimarySlot1    int `db:"rune_primary_slot1"`
	RunePrimarySlot2    int `db:"rune_primary_slot2"`
	RunePrimarySlot3    int `db:"rune_primary_slot3"`
	RuneSecondaryPath   int `db:"rune_secondary_path"`
	RuneSecondarySlot1  int `db:"rune_secondary_slot1"`
	RuneSecondarySlot2  int `db:"rune_secondary_slot2"`
	RuneShardSlot1      int `db:"rune_shard_slot1"`
	RuneShardSlot2      int `db:"rune_shard_slot2"`
	RuneShardSlot3      int `db:"rune_shard_slot3"`

	Kills       int `db:"kills"`
	Deaths      int `db:"deaths"`
	Assists     int `db:"assists"`
	CreepScore  int `db:"creep_score"`
	VisionScore int `db:"vision_score"`
	GoldEarned  int `db:"gold_earned"`
	GoldSpent   int `db:"gold_spent"`

	PlayerPosition string `db:"player_position"`
	ChampionLevel  int    `db:"champion_level"`
	ChampionID     int    `db:"champion_id"`
	ChampionName   string `db:"champion_name"`

	TotalDamageDealtToChampions int `db:"total_damage_dealt_to_champions"`
}

func (m MatchSummonerPostGame) GetKdaRatio() string {
	if m.Deaths == 0 {
		return "Perfect"
	}

	return fmt.Sprintf("%.2f", float32((m.Kills+m.Assists)/m.Deaths))
}

type MatchSummonerPostGameList []*MatchSummonerPostGame

type RiotMatchId string

func (id RiotMatchId) String() string {
	return string(id)
}

func (db *DB) ensureMatchlist(ctx context.Context, puuid string, start, count int) error {
	ids, err := db.riot.GetMatchHistory(ctx, puuid, start, count)
	if err != nil {
		return fmt.Errorf("riot match list: %w", err)
	}

	batch := &pgx.Batch{}
	for _, id := range ids {
		var found bool
		err := db.pool.QueryRow(ctx, "select 1 from match_info_records WHERE match_id = $1", id).Scan(&found)
		if err != nil {
			return fmt.Errorf("sql select: %w", err)
		}

		if found {
			logging.FromContext(ctx).Sugar().Debugf("not getting %s: already found", id)
			continue
		}

		match, err := db.riot.GetMatch(ctx, id)
		if err != nil {
			return fmt.Errorf("riot match: %w", err)
		}

		batchInsertRiotMatch(batch, match)
	}

	batchRes := db.pool.SendBatch(ctx, batch)
	defer batchRes.Close()

	for range batch.Len() {
		_, err := batchRes.Exec()
		if err != nil {
			return err
		}
	}

	return nil
}

func batchInsertRiotMatch(batch *pgx.Batch, m *riot.Match) {
	matchInfoRow := map[string]any{
		"match_id":      m.Metadata.MatchId,
		"game_date":     time.Unix(m.Info.GameStartTimestamp/1000, 0),
		"game_duration": time.Duration(m.Info.GameDuration) * time.Second,
		"game_patch":    m.Info.GameVersion,
	}

	batchInsertRow(batch, "match_info_records", matchInfoRow)

	for _, p := range m.Info.Participants {
		matchParticipantRow := map[string]any{
			"assists":              p.Assists,
			"champion_id":          p.ChampLevel,
			"champion_level":       p.ChampLevel,
			"creep_score":          p.TotalMinionsKilled,
			"deaths":               p.Deaths,
			"gold_earned":          p.GoldEarned,
			"gold_spent":           p.GoldSpent,
			"item0_id":             p.Item0,
			"item1_id":             p.Item1,
			"item2_id":             p.Item2,
			"item3_id":             p.Item3,
			"item4_id":             p.Item4,
			"item5_id":             p.Item5,
			"item6_id":             p.Item6,
			"kills":                p.Kills,
			"match_id":             m.Metadata.MatchId,
			"player_position":      p.Role,
			"player_win":           p.Win,
			"puuid":                p.PUUID,
			"rune_main_keystone":   p.Perks.Styles[0].Selections[0].Perk,
			"rune_main_path":       p.Perks.Styles[0].Style,
			"rune_main_slot1":      p.Perks.Styles[0].Selections[1].Perk,
			"rune_main_slot2":      p.Perks.Styles[0].Selections[2].Perk,
			"rune_main_slot3":      p.Perks.Styles[0].Selections[3].Perk, // THESE ARE ALL TODO
			"rune_secondary_path":  p.Perks.Styles[1].Style,
			"rune_secondary_slot1": p.Perks.Styles[0].Selections[1].Perk,
			"rune_secondary_slot2": p.Perks.Styles[0].Selections[1].Perk,
			"rune_shard_slot1":     p.Perks.StatPerks.Offense,
			"rune_shard_slot2":     p.Perks.StatPerks.Flex,
			"rune_shard_slot3":     p.Perks.StatPerks.Defense,
			"spell0_id":            p.Summoner1ID,
			"spell1_id":            p.Summoner2ID,
			"team_id":              p.TeamID,
			"vision_score":         p.VisionScore,
		}

		batchInsertRow(batch, "match_participant_records", matchParticipantRow)
	}

	for _, t := range m.Info.Teams {
		matchTeamRow := map[string]any{
			"match_id":               m.Metadata.MatchId,
			"team_id":                t.TeamId,
			"team_win":               t.Win,
			"team_surrendered":       false, // TODO
			"team_early_surrendered": false,
		}

		batchInsertRow(batch, "match_team_records", matchTeamRow)
	}
}
