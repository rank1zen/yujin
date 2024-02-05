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

	db := postgresql.NewQuery(pool)

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
		id, err := db.SummonerV4.InsertSummonerRecord(ctx, &test.arg, test.ts)
		assert.NoError(t, err)

		record, err := db.SummonerV4.SelectSummonerRecord(ctx, id)
		if assert.NoError(t, err) {
			assert.Equal(t, id, record.RecordId)
			assert.Equal(t, test.ts, record.RecordDate)
		}
	}
}

func TestGetSummonerRecord(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	addr := postgresql.NewDockerResource(t)

	pool, err := postgresql.BackoffRetryPool(ctx, addr)
	require.NoError(t, err)

	db := postgresql.NewQuery(pool)

	err = postgresql.Migrate(ctx, pool)
	require.NoError(t, err)

	ts := time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC)
	for i:= 1 ; i <= 100; i++ {
		db.SummonerV4.InsertSummonerRecord(ctx, &postgresql.SummonerRecordArg{Name: "Pobelter"}, ts)
	}

	records, err := db.SummonerV4.SelectSummonerRecordsByName(ctx, "Pobelter")
	assert.NoError(t, err)

	newest, err := db.SummonerV4.SelectSummonerRecordNewestByName(ctx, "Pobelter")
	assert.NoError(t, err)

	t.Log(records, newest)
}
