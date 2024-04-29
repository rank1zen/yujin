package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// Bulk insert a list of records, returns how many were inserted
func insertMatchParticipantRecords(db pgDB) func(context.Context, []*MatchParticipantRecord) (int64, error) {
	return func(ctx context.Context, records []*MatchParticipantRecord) (int64, error) {
		// turn the list of records into a 2d list of values
		var rows [][]any
		for _, r := range records {
			rows = append(rows, []any{
				r.MatchId, r.Puuid, r.ParticipantId, r.TeamId, r.SummonerName,
				r.SummonerLevel, r.Position, r.ChampId, r.ChampName, r.ChampLevel,
				r.Kills, r.Deaths, r.Assists, r.CreepScore, r.GoldEarned,
				r.VisionScore, r.WardsPlaced, r.ControlWardsPlaced, r.FirstBloodAssist,
				r.FirstTowerAssist, r.TurretTakeDowns, r.PhysicalDamageDealtToChampions,
				r.MagicDamageDealtToChampions, r.TrueDamageDealtToChampions,
				r.TotalDamageDealtToChampions, r.TotalDamageTaken, r.TotalHealsOnTeammates,
			})
		}

		count, err := db.CopyFrom(ctx,
			pgx.Identifier{"matchparticipantrecords"},
			[]string{
				"match_id", "puuid", "participant_id", "team_id", "summoner_name",
				"summoner_level", "position", "champion_id", "champion_name", "champion_level",
				"kills", "deaths", "assists", "creep_score", "gold_earned",
				"vision_score", "wards_placed", "control_wards_placed", "first_blood_assist",
				"first_tower_assist", "turret_takedowns", "physical_damage_dealt_to_champions",
				"magic_damage_dealt_to_champions", "true_damage_dealt_to_champions",
				"total_damage_dealt_to_champions", "total_damage_taken", "total_heals_on_teammates",
			},
			pgx.CopyFromRows(rows))
		if err != nil {
			return 0, fmt.Errorf("insert match participant: %w", err)
		}

		return count, err
	}
}

// GetMatchlist returns match IDs associated with a puuid
func getMatchlist(db pgDB) func(context.Context, string) ([]string, error) {
	return func(ctx context.Context, puuid string) ([]string, error) {
		rows, _ := db.Query(ctx, `
                        SELECT match_id
                        FROM MatchParticipantRecords
                        WHERE puuid = $1
                `, puuid)

		defer rows.Close()
		records, err := pgx.CollectRows(rows, pgx.RowToStructByPos[string])
		if err != nil {
			return nil, fmt.Errorf("get matchlist: %w", err)
		}

		return records, nil
	}
}
