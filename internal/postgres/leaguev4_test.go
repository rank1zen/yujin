package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/yujin/internal/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertSoloqRecord(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	addr := NewDockerResource(t)

	pool, err := postgres.NewBackoffPool(ctx, addr)
	require.NoError(t, err)

	db := postgres.NewQuery(pool)

	require.NoError(t, err)

	tests := []struct {
		arg postgres.SoloqRecordArg
		ts  time.Time
	}{
		{
			arg: postgres.SoloqRecordArg{},
			ts:  time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC),
		},
		{
			arg: postgres.SoloqRecordArg{},
			ts:  time.Now(),
		},
	}

	for _, test := range tests {
		id, err := db.LeagueV4.InsertSoloqRecord(ctx, &test.arg, test.ts)
		assert.NoError(t, err)

		record, err := db.LeagueV4.SelectSoloqRecord(ctx, id)
		if assert.NoError(t, err) {
			assert.Equal(t, id, record.RecordId)
			assert.Equal(t, test.ts, record.RecordDate)
		}
	}
}
