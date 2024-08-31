package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/internal/logging"
	"github.com/rank1zen/yujin/internal/riot"
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

type MatchSummonerPostGameList []*MatchSummonerPostGame

type RiotMatchId string

func (id RiotMatchId) String() string {
	return string(id)
}

func (db *DB) ensureMatchlist(ctx context.Context, puuid string, start, count int) error {
	ids, err := db.riot.GetMatchIdsByPuuid(ctx, puuid, start, count)
	if err != nil {
		return fmt.Errorf("riot match list: %w", err)
	}

	batch := &pgx.Batch{}
	for _, id := range ids {
		var found bool
		err := db.pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM
				match_info_records
			WHERE
				match_id = $1
		);
		`, id).Scan(&found)
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

// NOTE: A lot of rows are missing
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
			"item0_id": p.Item0,
			"item1_id": p.Item1,
			"item2_id": p.Item2,
			"item3_id": p.Item3,
			"item4_id": p.Item4,
			"item5_id": p.Item5,
			"item6_id": p.Item6,

			"spell1_id": p.Summoner1ID,
			"spell2_id": p.Summoner2ID,

			"kills":        p.Kills,
			"assists":      p.Assists,
			"deaths":       p.Deaths,
			"creep_score":  p.TotalMinionsKilled,
			"vision_score": p.VisionScore,

			"match_id":        m.Metadata.MatchId,
			"puuid":           p.PUUID,
			"team_id":         p.TeamID,
			"participant_id":  p.ParticipantID,
			"player_position": p.Role,
			"player_win":      p.Win,
			"champion_level":  p.ChampLevel,
			"champion_id":     p.ChampionID,
			"champion_name":   p.ChampionName,

			"rune_primary_path":     p.Perks.Styles[0].Style,
			"rune_primary_keystone": p.Perks.Styles[0].Selections[0].Perk,
			"rune_primary_slot1":    p.Perks.Styles[0].Selections[1].Perk,
			"rune_primary_slot2":    p.Perks.Styles[0].Selections[2].Perk,
			"rune_primary_slot3":    p.Perks.Styles[0].Selections[3].Perk,
			"rune_secondary_path":   p.Perks.Styles[1].Style,
			"rune_secondary_slot1":  p.Perks.Styles[1].Selections[0].Perk,
			"rune_secondary_slot2":  p.Perks.Styles[1].Selections[1].Perk,
			"rune_shard_slot1":      p.Perks.StatPerks.Offense,
			"rune_shard_slot2":      p.Perks.StatPerks.Flex,
			"rune_shard_slot3":      p.Perks.StatPerks.Defense,

			"physical_damage_dealt":              p.PhysicalDamageDealt,
			"physical_damage_dealt_to_champions": p.PhysicalDamageDealtToChampions,
			"physical_damage_taken":              p.PhysicalDamageTaken,
			"magic_damage_dealt":                 p.MagicDamageDealt,
			"magic_damage_dealt_to_champions":    p.MagicDamageDealtToChampions,
			"magic_damage_taken":                 p.MagicDamageTaken,
			"true_damage_dealt":                  p.TrueDamageDealt,
			"true_damage_dealt_to_champions":     p.TrueDamageDealtToChampions,
			"true_damage_taken":                  p.TrueDamageTaken,
			"total_damage_dealt":                 p.TotalDamageDealt,
			"total_damage_dealt_to_champions":    p.TotalDamageDealtToChampions,
			"total_damage_taken":                 p.TotalDamageTaken,
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
