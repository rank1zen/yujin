package database

import (
	"context"

	"github.com/rank1zen/yujin/pkg/riot"
)

// will return nil if not found
func findSoloqRank(entries riot.LeagueEntryList) *riot.LeagueEntry {
	for _, entry := range entries {
		if entry.QueueType == "a" {
			return entry
		}
	}

	return nil
}

func (db *DB) updateSummonerRankRecord(ctx context.Context, m *riot.LeagueEntry) error {
	row := map[string]any{
		"summoner_id":             m.SummonerId,
		"league_id":     m.LeagueId,
		"tier":          m.Tier,
		"division":      m.Rank,
		"league_points": m.LeaguePoints,
		"number_wins":   m.Wins,
		"number_losses": m.Losses,
	}

	return queryInsertRow(ctx, db.pool, "league_records", row)
}
