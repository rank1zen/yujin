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

func (db *DB) UpdateSummonerRankRecord(ctx context.Context, summonerID string) error {
	entries, err := db.riot.GetLeagueEntriesForSummoner(ctx, summonerID)
	if err != nil {
		return err
	}

	row := map[string]any{"summoner_id": summonerID}

	soloq := findSoloqRank(entries)
	if soloq != nil {
		row["league_id"] = soloq.LeagueId
		row["tier"] = soloq.Tier
		row["division"] = soloq.Rank
		row["league_points"] = soloq.LeaguePoints
		row["number_wins"] = soloq.Wins
		row["number_losses"] = soloq.Losses
	}

	_, err = insertRow(ctx, db.pool, "league_records", row)
	if err != nil {
		return err
	}

	return nil
}
