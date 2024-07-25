package database

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/pkg/riot"
)

// CheckExistingMatchRecords returns all existing match ids found in DB
func (db *DB) CheckExistingMatchRecords(ctx context.Context, puuid string) ([]string, error) {
	return Select(
		ctx,
		db.pool,
		`SELECT match_id FROM match_info_records WHERE puuid = $1`,
		[]any{puuid},
		pgx.RowToStructByName[string],
	)
}

func (db *DB) FetchMatchRecords(ctx context.Context, puuid string, matchIDs []string) error {
	var batch pgx.Batch

	existing, err := db.CheckExistingMatchRecords(ctx, puuid)
	if err != nil {
		return err
	}

	for _, id := range matchIDs {
		if !slices.Contains(existing, id) {
			matchDto, err := db.riot.GetMatch(ctx, id)
			if err != nil {
				return err
			}
			batchMatch(&batch, matchDto)
		}
	}

	res := db.pool.SendBatch(ctx, &batch)
	defer res.Close()

	for range batch.Len() {
		_, err := res.Exec()
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) CreateMatches(ctx context.Context, m riot.MatchDtoList) error {
	var batch pgx.Batch

	for _, match := range m {
		batchMatch(&batch, match)
	}

	br := db.pool.SendBatch(ctx, &batch)
	defer br.Close()

	for range batch.Len() {
		_, err := br.Exec()
		if err != nil {
			return err
		}
	}

	return nil
}

// getMatchListPlayer fetches each match corresponding to each matchIDs.
// Calls the Riot API and inserts matches as needed.
// NOTE: This is currently not in transaction
func (db *DB) getMatchListPlayer(ctx context.Context, puuid string, matchIDs []string) ([]*ProfileMatch, error) {
	matches := make([]*ProfileMatch, len(matchIDs))

	for i, id := range matchIDs {
		match, err := db.getMatchPlayer(ctx, puuid, id)
		if err != nil {
			return nil, fmt.Errorf("failed %s: %w", id, err)
		}
		matches[i] = match
	}

	return matches, nil
}

func batchMatch(batch *pgx.Batch, m *riot.MatchDto) {
	matchID := m.Metadata.MatchId

	var sql string
	var args []any

	sql, args = matchInfoQuery(m)
	batch.Queue(sql, args...)

	for _, p := range m.Info.Participants {
		sql, args = matchParticipantQuery(matchID, p)
		batch.Queue(sql, args...)
		sql, args = matchItemQuery(matchID, p)
		batch.Queue(sql, args...)
		sql, args = matchSummonerSpellQuery(matchID, p)
		batch.Queue(sql, args...)
	}

	for _, t := range m.Info.Teams {
		sql, args = matchTeamQuery(matchID, t)
		batch.Queue(sql, args...)
		sql, args = matchObjectiveQuery(matchID, t)
		batch.Queue(sql, args...)

		for _, ban := range t.Bans {
			sql, args = matchBanQuery(matchID, t.TeamId, ban)
			batch.Queue(sql, args...)
		}
	}
}

func (db *DB) getMatchPlayer(ctx context.Context, puuid, matchID string) (*ProfileMatch, error) {
	batch := new(pgx.Batch)

	sql := `
	SELECT
		match_id, game_date, game_duration, game_patch,
		player_win, player_position, kills, deaths, assists, creep_score, champion_level, champion_id, vision_score,
		items_arr, spells_arr, runes_arr
	FROM match_participant_simple
	WHERE match_id = $1 and puuid = $2;
	`

	row, _ := db.pool.Query(ctx, sql, matchID, puuid)
	match, err := pgx.CollectExactlyOneRow(row, pgx.RowToAddrOfStructByNameLax[ProfileMatch])

	var riotMatch *riot.MatchDto

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		riotMatch, err = db.riot.GetMatch(ctx, matchID)
		if err != nil {
			return nil, fmt.Errorf("failed %s: %w", matchID, err)
		}
	case err != nil:
		return nil, fmt.Errorf("failed to check db: %w", err)
	default:
		return match, nil
	}

	batchMatch(batch, riotMatch)
	batch.Queue(sql, matchID, puuid)

	batchRes := db.pool.SendBatch(ctx, batch)
	defer batchRes.Close()

	for range batch.Len() - 1 {
		tag, err := batchRes.Exec()
		if err != nil {
			return nil, fmt.Errorf("batch insert: %v %w", tag, err)
		}
	}

	row, _ = batchRes.Query()
	match, err = pgx.CollectExactlyOneRow(row, pgx.RowToAddrOfStructByNameLax[ProfileMatch])
	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}

	return match, nil
}
