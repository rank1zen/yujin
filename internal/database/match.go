package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/internal/logging"
	"github.com/rank1zen/yujin/internal/pgxutil"
	"github.com/rank1zen/yujin/internal/riot"
)

func ensureMatchList(ctx context.Context, db pgxutil.Conn, r *riot.Client, puuid string, start, count int) error {
	ids, err := r.GetMatchIdsByPuuid(ctx, puuid, start, count)
	if err != nil {
		return fmt.Errorf("fetching riot: %w", err)
	}

	batch := &pgx.Batch{}
	for _, id := range ids {
		var found bool
		err := db.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM match_info_records WHERE match_id = $1);
		`, id).Scan(&found)
		if err != nil {
			return fmt.Errorf("sql select: %w", err)
		}

		if found {
			logging.FromContext(ctx).Sugar().Debugf("not getting %s: already found", id)
			continue
		}

		match, err := r.GetMatch(ctx, id)
		if err != nil {
			return fmt.Errorf("riot match: %w", err)
		}

		batchInsertRiotMatch(batch, match)
	}

	batchRes := db.SendBatch(ctx, batch)
	defer batchRes.Close()

	for range batch.Len() {
		_, err := batchRes.Exec()
		if err != nil {
			return fmt.Errorf("inserting match: %w", err)
		}
	}

	return nil
}

func batchInsertRiotMatch(batch *pgx.Batch, m *riot.Match) {
	matchInfoRow := map[string]any{
		"match_id":      m.Metadata.MatchId,
		"data_version":  m.Metadata.DataVersion,
		"game_date":     time.Unix(m.Info.GameEndTimestamp/1000, 0),
		"game_duration": time.Duration(m.Info.GameDuration) * time.Second,
		"game_patch":    m.Info.GameVersion,
	}

	pgxutil.BatchInsertRow(batch, "match_info_records", matchInfoRow)

	for _, p := range m.Info.Participants {
		matchParticipantRow := map[string]any{
			"match_id":       m.Metadata.MatchId,
			"puuid":          p.PUUID,
			"team_id":        p.TeamID,
			"participant_id": p.ParticipantID,

			"player_position": p.Role,
			"champion_level":  p.ChampLevel,
			"champion_id":     p.ChampionID,
			"champion_name":   p.ChampionName,

			"kills":        p.Kills,
			"assists":      p.Assists,
			"deaths":       p.Deaths,
			"creep_score":  p.TotalMinionsKilled + p.NeutralMinionsKilled,
			"vision_score": p.VisionScore,

			"spell1_id": p.Summoner1ID,
			"spell2_id": p.Summoner2ID,

			"item0_id": p.Item0,
			"item1_id": p.Item1,
			"item2_id": p.Item2,
			"item3_id": p.Item3,
			"item4_id": p.Item4,
			"item5_id": p.Item5,
			"item6_id": p.Item6,

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

		pgxutil.BatchInsertRow(batch, "match_participant_records", matchParticipantRow)
	}

	for _, t := range m.Info.Teams {
		matchTeamRow := map[string]any{
			"match_id": m.Metadata.MatchId,
			"team_id":  t.TeamId,
			"win": t.Win,
		}

		pgxutil.BatchInsertRow(batch, "match_team_records", matchTeamRow)
	}
}
