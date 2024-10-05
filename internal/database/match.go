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

type MatchParticipant struct {
	Name       string
	Rank       string
	Kills      string
	Deaths     string
	Assists    string
	CreepScore string
	CsPer10    string
	Damage     string

	ChampionIcon      string
	RunePrimaryIcon   string
	RuneSecondaryIcon string
	Spell1Icon        string
	Spell2Icon        string
	ItemIcons         []*string
}

type Match struct {
	MatchId      string
	GamePatch    string
	GameDuration string
	GameDate     string
	RedSide      []MatchParticipant
	BlueSide     []MatchParticipant
}

func (db *DB) MatchGet(ctx context.Context, name, matchID string) (Match, error) {
	var m Match

	batch := &pgx.Batch{}

	batch.Queue(`
	SELECT
		match_id,
		game_patch,
		EXTRACT(MINUTE FROM game_duration) || 'm ' || EXTRACT(SECOND FROM game_duration) || 's' AS game_duration,
		TO_CHAR(game_date, 'MM-DD HH24:MI') AS game_date
	FROM
		match_info_records
	WHERE
		match_id = $1;
	`, matchID).QueryRow(func(row pgx.Row) error {
		err := row.Scan(&m.MatchId, &m.GamePatch, &m.GameDuration, &m.GameDate)
		if err != nil {
			return fmt.Errorf("getting match info: %w", err)
		}
		return nil
	})

	getPostGame := func(teamID int, dst *[]MatchParticipant) {
		batch.Queue(`
		SELECT
			participant_name AS name,
			'???' AS rank,
			kills, deaths, assists,
			creep_score, TO_CHAR(60 * creep_score / EXTRACT(epoch FROM game_duration), 'FM99999.0') AS cs_per_10,
			total_damage_dealt_to_champions AS damage,
			FORMAT('https://cdn.communitydragon.org/14.16.1/champion/%s/square', champion_id) AS champion_icon_url,
			array[item0_id, item1_id, item2_id, item3_id, item4_id, item5_id] as items,
			array[spell1_id, spell2_id] as spells,
			rune_primary_keystone AS rune_primary,
			rune_secondary_path AS rune_secondary
		FROM
			profile_matches
		WHERE 1=1
			AND match_id = $1
			AND team_id = $2
		ORDER BY
			participant_id;
		`, matchID, teamID).Query(func(rows pgx.Rows) error {
			collectedRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[MatchParticipant])
			if err != nil {
				return fmt.Errorf("getting team: %w", err)
			}

			if len(collectedRows) != 5 {
				logging.FromContext(ctx).DPanic("team does not have 5")
			}

			*dst = collectedRows
			return nil
		})
	}

	getPostGame(100, &m.BlueSide)
	getPostGame(200, &m.RedSide)

	err := db.pool.SendBatch(ctx, batch).Close()
	if err != nil {
		return Match{}, err
	}

	return m, nil
}

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

		matchBatchInsert(batch, match)
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

type riotMatchRows struct {
	MatchInfoRow         map[string]any
	MatchParticipantRows []map[string]any
	MatchTeamRows        []map[string]any
}

func riotMatchToRows(m *riot.Match) riotMatchRows {
	var matchParticipantRows []map[string]any
	for _, p := range m.Info.Participants {
		row := map[string]any{
			"id":       p.ParticipantID,
			"match_id": m.Metadata.MatchId,
			"team_id":  p.TeamID,
			"puuid":    p.PUUID,
			"name":     p.RiotIDGameName + "#" + p.RiotIDTagline,

			"kills":           p.Kills,
			"assists":         p.Assists,
			"deaths":          p.Deaths,
			"creep_score":     p.TotalMinionsKilled + p.NeutralMinionsKilled,
			"vision_score":    p.VisionScore,
			"gold_earned":     p.GoldEarned,
			"gold_spent":      p.GoldSpent,
			"player_position": p.Role,
			"champion_level":  p.ChampLevel,
			"champion_id":     p.ChampionID,
			"champion_name":   p.ChampionName,

			"spell1_id": p.Summoner1ID,
			"spell2_id": p.Summoner2ID,

			"items": []int{p.Item0, p.Item1, p.Item3, p.Item3, p.Item4, p.Item5, p.Item6},

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

		matchParticipantRows = append(matchParticipantRows, row)
	}

	var matchTeamRows []map[string]any
	for _, t := range m.Info.Teams {
		row := map[string]any{
			"id":       t.TeamId,
			"match_id": m.Metadata.MatchId,
			"win":      t.Win,
		}

		matchTeamRows = append(matchTeamRows, row)
	}

	return riotMatchRows{
		MatchInfoRow: map[string]any{
			"id":            m.Metadata.MatchId,
			"data_version":  m.Metadata.DataVersion,
			"game_date":     riotUnixToDate(m.Info.GameEndTimestamp),
			"game_duration": time.Duration(m.Info.GameDuration) * time.Second,
			"game_patch":    m.Info.GameVersion,
		},
		MatchParticipantRows: matchParticipantRows,
		MatchTeamRows:        matchTeamRows,
	}
}

func matchInsert(ctx context.Context, db pgxutil.Conn, m *riot.Match) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	rows := riotMatchToRows(m)

	err = pgxutil.QueryInsertRow(ctx, db, "matches", rows.MatchInfoRow)
	if err != nil {
		return err
	}

	for _, p := range rows.MatchParticipantRows {
		err := pgxutil.QueryInsertRow(ctx, db, "match_participants", p)
		if err != nil {
			return err
		}
	}

	for _, p := range rows.MatchTeamRows {
		err := pgxutil.QueryInsertRow(ctx, db, "match_teams", p)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func matchBatchInsert(batch *pgx.Batch, m *riot.Match) {
	rows := riotMatchToRows(m)

	pgxutil.BatchInsertRow(batch, "matches", rows.MatchInfoRow)

	for _, p := range rows.MatchParticipantRows {
		pgxutil.BatchInsertRow(batch, "match_participants", p)
	}

	for _, p := range rows.MatchTeamRows {
		pgxutil.BatchInsertRow(batch, "match_teams", p)
	}
}
