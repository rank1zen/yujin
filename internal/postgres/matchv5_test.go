package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/yujin/internal/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertMatchRecord(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	addr := NewDockerResource(t)
	pool := MustConnectionPool(t, addr)
	q := postgres.NewMatchV5Query(pool)

	tests := []struct {
		test string
		arg  postgres.MatchRecordArg
		err  bool
	}{
		{
			test: "Empty Arg",
			arg:  postgres.MatchRecordArg{},
			err:  false,
		},
		{
			test: "Standard Arg",
			arg: postgres.MatchRecordArg{
				RecordDate: time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC),
				MatchId:    "10",
				StartDate:  time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC),
				Duration:   time.Duration(25 * time.Minute),
				Surrender:  true,
				Patch:      "20",
			},
			err: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.test, func(t *testing.T) {
			_, err := q.InsertMatch(ctx, &tc.arg)
			assert.NoError(t, err)
		})
	}
}

func TestInsertMatchTeamRecord(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	addr := NewDockerResource(t)

	pool := MustConnectionPool(t, addr)

	q := postgres.NewMatchV5Query(pool)

	_, err := q.InsertMatch(ctx, &postgres.MatchRecordArg{})
	require.NoError(t, err)

	tests := []struct {
		arg postgres.MatchTeamRecordArg
	}{
		{
			arg: postgres.MatchTeamRecordArg{
				Objective: []*postgres.TeamObjective{
					{Name: "hi"},
					{Name: "bye"},
				},
				Bans: []*postgres.TeamBan{
					{},
				},
			},
		},
	}

	for _, tc := range tests {
		id, err := q.InsertMatchTeam(ctx, &tc.arg)
		if assert.NoError(t, err) {
			t.Log(id)
		}

		record, err := q.SelectMatchTeam(ctx, id)
		if assert.NoError(t, err) {
			t.Log(record[0].Bans)
			t.Log(record[0].Objectives)
		}
	}
}
