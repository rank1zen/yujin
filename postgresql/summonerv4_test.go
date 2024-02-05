package postgresql_test

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/yujin/postgresql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertSummonerRecord(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	addr := postgresql.NewDockerResource(t)

	pool, err := postgresql.BackoffRetryPool(ctx, addr)
	require.NoError(t, err)

	db := postgresql.NewQueries(pool)

	err = postgresql.Migrate(ctx, pool)
	require.NoError(t, err)

	tests := []struct {
		arg postgresql.SummonerRecordArg
		ts  time.Time
	}{
		{
			arg: postgresql.SummonerRecordArg{},
			ts:  time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC),
		},
		{
			arg: postgresql.SummonerRecordArg{},
			ts:  time.Now(),
		},
	}

	for _, test := range tests {
		id, err := db.InsertSummonerRecord(ctx, &test.arg, test.ts)
		assert.NoError(t, err)

		record, err := db.SelectSummonerRecord(ctx, id)
		if assert.NoError(t, err) {
			assert.Equal(t, record.RecordId, id)
			assert.Equal(t, record.RecordDate, test.ts)
		}
	}

	records, err := db.SelectSummonerRecordsByName(ctx, "")
	assert.NoError(t, err)
	t.Log(records)
}
