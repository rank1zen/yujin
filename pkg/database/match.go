package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/pkg/riot"
)

type RiotMatchId string

func (id RiotMatchId) String() string {
	return string(id)
}

type RiotMatchIdList []RiotMatchId

func (db *DB) ensureMatchlist(ctx context.Context, puuid RiotPuuid, start, count int) error {
	ids, err := db.riot.GetMatchHistory(ctx, puuid.String(), start, count)
	if err != nil {
		return err
	}

	batch := &pgx.Batch{}
	for _, id := range ids {
		match, err := db.riot.GetMatch(ctx, id)
		if err != nil {
			return err
		}

		batchMatch(batch, match)
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

func batchMatch(batch *pgx.Batch, m *riot.MatchDto) {
	batchInsertRow(batch, "match_info_records", map[string]any{
		"match_id":      m.Metadata.MatchId,
		"game_date":     time.Unix(m.Info.GameStartTimestamp/1000, 0),
		"game_duration": time.Duration(m.Info.GameDuration) * time.Second,
		"game_patch":    m.Info.GameVersion,
	})

	for _, p := range m.Info.Participants {
		batchInsertRow(batch, "match_participant_records", map[string]any{
			"match_id":             m.Metadata.MatchId,
			"puuid":                p.PUUID,
			"team_id":              p.TeamID,
			"player_win":           p.Win,
			"player_position":      p.Role,
			"kills":                p.Kills,
			"deaths":               p.Deaths,
			"assists":              p.Assists,
			"creep_score":          p.TotalMinionsKilled,
			"vision_score":         p.VisionScore,
			"gold_earned":          p.GoldEarned,
			"champion_level":       p.ChampLevel,
			"champion_id":          p.ChampLevel,
			"rune_main_path":       p.Perks.Styles[0].Style,
			"rune_main_keystone":   p.Perks.Styles[0].Selections[0].Perk,
			"rune_main_slot1":      p.Perks.Styles[0].Selections[1].Perk,
			"rune_main_slot2":      p.Perks.Styles[0].Selections[2].Perk,
			"rune_main_slot3":      p.Perks.Styles[0].Selections[3].Perk, // THESE ARE ALL TODO
			"rune_secondary_path":  p.Perks.Styles[1].Style,
			"rune_secondary_slot1": p.Perks.Styles[0].Selections[1].Perk,
			"rune_secondary_slot2": p.Perks.Styles[0].Selections[1].Perk,
			"rune_shard_slot1":     p.Perks.StatPerks.Offense,
			"rune_shard_slot2":     p.Perks.StatPerks.Flex,
			"rune_shard_slot3":     p.Perks.StatPerks.Defense,
		})

		spell := func(batch *pgx.Batch, id, slot, casts int) {
			row := map[string]any{}

			row["match_id"] = m.Metadata.MatchId
			row["puuid"] = p.PUUID

			row["spell_slot"] = slot
			row["spell_id"] = id
			row["spell_casts"] = casts

			batchInsertRow(batch, "match_summonerspell_records", row)
		}

		spell(batch, p.Summoner1ID, 1, p.Summoner1Casts)
		spell(batch, p.Summoner2ID, 2, p.Summoner2Casts)

		item := func(batch *pgx.Batch, id, slot int) {
			batchInsertRow(batch, "match_item_records", map[string]any{
				"match_id":  m.Metadata.MatchId,
				"puuid":     p.PUUID,
				"item_id":   id,
				"item_slot": slot,
			})
		}

		item(batch, p.Item0, 0)
		item(batch, p.Item1, 1)
		item(batch, p.Item2, 2)
		item(batch, p.Item3, 3)
		item(batch, p.Item4, 4)
		item(batch, p.Item5, 5)
		item(batch, p.Item6, 6)
	}

	for _, t := range m.Info.Teams {
		batchInsertRow(batch, "match_team_records", map[string]any{
			"match_id":               m.Metadata.MatchId,
			"team_id":                t.TeamId,
			"team_win":               t.Win,
			"team_surrendered":       false, // TODO
			"team_early_surrendered": false,
		})

		obj := func(batch *pgx.Batch, name, first, kills any) {
			batchInsertRow(batch, "match_objective_records", map[string]any{
				"match_id": m.Metadata.MatchId,
				"team_id":  t.TeamId,
				"name":     name,
				"first":    first,
				"kills":    kills,
			})
		}

		obj(batch, "Baron", t.Objectives.Baron.First, t.Objectives.Baron.Kills)
		obj(batch, "RiftHerald", t.Objectives.RiftHerald.First, t.Objectives.RiftHerald.Kills)
		obj(batch, "Dragon", t.Objectives.Dragon.First, t.Objectives.Dragon.Kills)
		obj(batch, "Inhibitor", t.Objectives.Inhibitor.First, t.Objectives.Inhibitor.Kills)
		obj(batch, "Tower", t.Objectives.Tower.First, t.Objectives.Tower.Kills)
		obj(batch, "Champion", t.Objectives.Champion.First, t.Objectives.Champion.Kills)

		for _, ban := range t.Bans {
			batchInsertRow(batch, "match_ban_records", map[string]any{
				"match_id":    m.Metadata.MatchId,
				"team_id":     t.TeamId,
				"champion_id": ban.ChampionId,
				"turn":        ban.PickTurn,
			})
		}
	}
}
