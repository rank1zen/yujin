package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type matchV5Query struct {
	db pgDB
}

func NewMatchV5Query(db pgDB) MatchV5Query {
	return &matchV5Query{db: db}
}

func (q *matchV5Query) GetMatchRecordsMatchingIds(ctx context.Context, ids []string) ([]*MatchRecord, []string, error) {
	return getMatchRecordsMatchingIds(q.db)(ctx, ids)
}

func (q *matchV5Query) GetMatchRecords(ctx context.Context, filters ...RecordFilter) ([]*MatchRecord, error) {
        return []*MatchRecord{}, nil
	//return getMatchRecords(q.db)(ctx, filters...)
}

func (q *matchV5Query) GetMatchlist(ctx context.Context, puuid string) ([]string, error) {
	return getMatchlist(q.db)(ctx, puuid)
}

func (q *matchV5Query) InsertFullMatchRecords(ctx context.Context, records []*FullMatchRecord) (int64, error) {
	return insertFullMatchRecords(q.db)(ctx, records)
}

func (q *matchV5Query) InsertMatchRecords(ctx context.Context, records []*MatchRecord) (int64, error) {
	return insertMatchRecords(q.db)(ctx, records)
}

func (q *matchV5Query) InsertMatchTeamRecords(ctx context.Context, records []*MatchTeamRecord) (int64, error) {
	return insertMatchTeamRecords(q.db)(ctx, records)
}

func (q *matchV5Query) InsertMatchBanRecords(ctx context.Context, records []*MatchBanRecord) (int64, error) {
	return insertMatchBanRecords(q.db)(ctx, records)
}

func (q *matchV5Query) InsertMatchParticipantRecords(ctx context.Context, records []*MatchParticipantRecord) (int64, error) {
	return insertMatchParticipantRecords(q.db)(ctx, records)
}

// Inserts a collection of records, we kinda assume that the match_id exists for each participant, team, etc...
func insertFullMatchRecords(db pgDB) func(context.Context, []*FullMatchRecord) (int64, error) {
	return func(ctx context.Context, records []*FullMatchRecord) (int64, error) {
		var matchs []*MatchRecord
		var teams []*MatchTeamRecord
		var bans []*MatchBanRecord
		var objs []*MatchObjectiveRecord
		var participants []*MatchParticipantRecord

		for _, r := range records {
			matchs = append(matchs, r.Metadata)
			teams = append(teams, r.Teams...)
			bans = append(bans, r.Bans...)
			objs = append(objs, r.Objectives...)
			participants = append(participants, r.Participants...)
		}

		count := int64(0)
		err := pgx.BeginFunc(ctx, db, func(tx pgx.Tx) error {
			c, err := insertMatchRecords(db)(ctx, matchs)
			if err != nil {
				return err
			}
			count += c

			return nil
		})

		if err != nil {
			return 0, fmt.Errorf("insert full match: %w", err)
		}

		// the idea is to insert all the rows associated with a match
		return count, nil
	}
}

// Return records with ids that are found, also return a list of ids not found
func getMatchRecordsMatchingIds(db pgDB) func(context.Context, []string) ([]*MatchRecord, []string, error) {
	return func(ctx context.Context, ids []string) ([]*MatchRecord, []string, error) {

		rows, _ := db.Query(ctx, `
                        SELECT
                                record_id, record_date, match_id, start_ts, duration, surrender, patch
                        FROM
                                MatchRecords
                        WHERE match_id = ANY($1)
                `, ids)

		defer rows.Close()
		records, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[MatchRecord])
		if err != nil {
			return nil, ids, fmt.Errorf("select match with ids: %w", err)
		}

		// ??? find the ids not found, there better be something else
		var remain []string
		for _, id := range ids {
			found := false
			for _, a := range records {
				if a.MatchId == id {
					found = true
				}
			}
			if !found {
				remain = append(remain, id)
			}
		}

		return records, remain, nil
	}
}
