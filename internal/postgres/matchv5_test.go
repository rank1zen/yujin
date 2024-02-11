package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/yujin/internal/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertMatchRecord(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	addr := postgres.NewDockerResource(t)

	pool := MustConnectionPool(t, addr)
	MustMigrate(t, pool)

	q := postgres.NewMatchV5Query(pool)

	tests := []struct {
		arg  postgres.MatchRecordArg
	}{
		{
			arg: postgres.MatchRecordArg{},
		},
	}

	for c, tc := range tests {
		id, err := q.InsertMatch(ctx, &tc.arg)
		assert.NoError(t, err, c)
		t.Log(id)
	}
}

func TestMatchTeamRecord(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	addr := postgres.NewDockerResource(t)

	pool := MustConnectionPool(t, addr)
	MustMigrate(t, pool)

	q := postgres.NewMatchV5Query(pool)

	_, err := q.InsertMatch(ctx, &postgres.MatchRecordArg{
	})
	require.NoError(t, err)

	tests := []struct {
		arg  postgres.MatchTeamRecordArg
		want string
	}{
		{
			arg: postgres.MatchTeamRecordArg{
				Objective: postgres.TeamObjective{},
				Bans:      postgres.TeamBan{},
			},
			want: "",
		},
	}

	for _, tc := range tests {
		id, err := q.InsertMatchTeam(ctx, &tc.arg)
		if assert.NoError(t, err) {
			assert.Equal(t, tc.want, id)
			t.Log(id)
		}
	}
}
