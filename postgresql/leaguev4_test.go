package postgresql_test

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/yujin/postgresql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertSoloqRecord(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	addr := postgresql.NewDockerResource(t)

	pool, err := postgresql.BackoffRetryPool(ctx, addr)
	require.NoError(t, err)

	db := postgresql.NewQuery(pool)

	err = postgresql.Migrate(ctx, pool)
	require.NoError(t, err)

	tests := []struct {
		arg postgresql.SoloqRecordArg
		ts  time.Time
	}{
		{
			arg: postgresql.SoloqRecordArg{},
			ts:  time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC),
		},
		{
			arg: postgresql.SoloqRecordArg{},
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
