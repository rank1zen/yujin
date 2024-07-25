package database

import (
	"fmt"
	"time"

	"github.com/rank1zen/yujin/pkg/riot"
)

func matchInfoQuery(m *riot.MatchDto) (string, []any) {
	// NOTE: When we get further along, we some fields will be missing and we have to adjust gameDuration
	// since that will be in millisecs
	// NOTE: the end time of a game should be the max time played of any player, check docs bro

	sql := `
	INSERT INTO match_info_records
		(match_id, game_date, game_duration, game_patch)
	VALUES
		($1, $2, $3, $4);
	`

	gameDate := time.Unix(m.Info.GameStartTimestamp/1000, 0)
	gameDuration := time.Duration(m.Info.GameDuration) * time.Second // p sure this is correct

	return sql, []any{
		m.Metadata.MatchId,
		gameDate,
		gameDuration,
		m.Info.GameVersion,
	}
}

func matchParticipantQuery(matchID string, m *riot.Participant) (string, []any) {
	sql := `
	INSERT INTO match_participant_records
		(match_id, puuid, player_win, player_position, kills, deaths, assists, creep_score,
		gold_earned, champion_level, champion_id)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
	`

	return sql, []any{
		matchID, m.PUUID, m.Win, m.Role, m.Kills,
		m.Deaths, m.Assists, m.TotalMinionsKilled, m.GoldEarned, m.ChampLevel,
		m.ChampionID,
	}
}

func matchRuneQuery(matchID, m *riot.Participant) (string, []any) {
	sql := `
	INSERT INTO match_rune_records
		(match_id, puuid, rune_id, rune_slot)
	VALUES
		($1, $2, $3, $4),
		($6, $7, $8, $9);
	`

	// 9 choices
	return sql, []any{
		matchID, m.PUUID,
	}
}

func matchSummonerSpellQuery(matchID string, m *riot.Participant) (string, []any) {
	sql := `
	INSERT INTO match_summonerspell_records
		(match_id, puuid, spell_slot, spell_id, spell_casts)
	VALUES
		($1, $2, $3, $4, $5),
		($6, $7, $8, $9, $10);
	`
	return sql, []any{
		matchID, m.PUUID, 1, m.Summoner1ID, m.Summoner1Casts,
		matchID, m.PUUID, 2, m.Summoner2ID, m.Summoner2Casts,
	}
}

func matchItemQuery(matchID string, m *riot.Participant) (string, []any) {
	sql := fmt.Sprintf(`
	INSERT INTO match_item_records
		(match_id, puuid, item_id, item_slot)
	VALUES
		($%d, $%d, $%d, $%d),
		($%d, $%d, $%d, $%d),
		($%d, $%d, $%d, $%d),
		($%d, $%d, $%d, $%d),
		($%d, $%d, $%d, $%d),
		($%d, $%d, $%d, $%d),
		($%d, $%d, $%d, $%d);
	`, bruh(28)...)

	return sql, []any{
		matchID, m.PUUID, m.Item0, 0,
		matchID, m.PUUID, m.Item1, 1,
		matchID, m.PUUID, m.Item2, 2,
		matchID, m.PUUID, m.Item3, 3,
		matchID, m.PUUID, m.Item4, 4,
		matchID, m.PUUID, m.Item5, 5,
		matchID, m.PUUID, m.Item6, 6,
	}
}

func matchTeamQuery(matchID string, m *riot.Team) (string, []any) {
	// FIXME: team surrender
	sql := `
	INSERT INTO match_team_records
		(match_id, team_id, team_win, team_surrendered, team_early_surrendered)
	VALUES
		($1, $2, $3, $4, $5)
	`
	return sql, []any{matchID, m.TeamId, m.Win, false, false}
}

func matchBanQuery(matchID string, teamID int, m *riot.TeamBan) (string, []any) {
	sql := `
	INSERT INTO match_ban_records
		(match_id, team_id, champion_id, turn)
	VALUES
		($1, $2, $3, $4)
	`
	return sql, []any{matchID, teamID, m.ChampionId, m.PickTurn}
}

func matchObjectiveQuery(matchID string, m *riot.Team) (string, []any) {
	sql := fmt.Sprintf(`
	INSERT INTO match_objective_records
		(match_id, team_id, name, first, kills)
	VALUES
		($%d, $%d, $%d, $%d, $%d),
		($%d, $%d, $%d, $%d, $%d),
		($%d, $%d, $%d, $%d, $%d),
		($%d, $%d, $%d, $%d, $%d),
		($%d, $%d, $%d, $%d, $%d),
		($%d, $%d, $%d, $%d, $%d);
	`, bruh(30)...)

	teamID := m.TeamId
	obj := m.Objectives

	return sql, []any{
		matchID, teamID, "Baron", obj.Baron.First, obj.Baron.Kills,
		matchID, teamID, "RiftHerald", obj.RiftHerald.First, obj.RiftHerald.Kills,
		matchID, teamID, "Dragon", obj.Dragon.First, obj.Dragon.Kills,
		matchID, teamID, "Inhibitor", obj.Inhibitor.First, obj.Inhibitor.Kills,
		matchID, teamID, "Tower", obj.Tower.First, obj.Tower.Kills,
		matchID, teamID, "Champion", obj.Champion.First, obj.Champion.Kills,
	}
}

func bruh(n int) []any {
	pos := make([]any, n)
	for i := 1; i <= n; i++ {
		pos[i-1] = i
	}
	return pos
}
