package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rank1zen/yujin/internal/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertSummonerRecord(t *testing.T) {
	t.Parallel()

	addr := postgres.NewDockerResource(t)

	pool, err := postgres.NewBackoffPool(context.Background(), addr)
	require.NoError(t, err)

	db := postgres.NewSummonerV4Query(pool)

	err = postgres.Migrate(context.Background(), pool)
	require.NoError(t, err)

	tests := []struct {
		tc   int
		arg  postgres.SummonerRecordArg
		date time.Time
	}{
		{
			tc:   1,
			arg:  postgres.SummonerRecordArg{},
			date: time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC),
		},
		{
			tc: 2,
			arg: postgres.SummonerRecordArg{Name: "testing"},
			date: time.Now(),
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("test case: %d", tc.tc), func(t *testing.T) {
			id, err := db.InsertSummonerRecord(context.Background(), &tc.arg, tc.date)
			assert.NoError(t, err)

			record, err := db.SelectSummonerRecord(context.Background(), id)
			if assert.NoError(t, err) {
				assert.Equal(t, id, record.RecordId)
				assert.Equal(t, tc.date, record.RecordDate)
				assert.Equal(t, tc.arg.Name, record.Name)
			}
		})
	}
}
