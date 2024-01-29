package postgresql_test

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/yujin/postgresql"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestPGX(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	addr := postgresql.NewDockerResource(t)

	log := zap.Must(zap.NewProduction())

	pool, err := postgresql.BackoffRetryPool(ctx, addr, log)
	assert.NoError(t, err)

	db := postgresql.NewQueries(pool)

	err = postgresql.Migrate(ctx, pool)
	assert.NoError(t, err)

	arg := &postgresql.SummonerRecord{
		RecordDate:    time.Now(),
		AccountId:     "hi",
		ProfileIconId: 120,
		RevisionDate:  1230,
		Name:          "ni",
		SummonerId:    "ad",
		Puuid:         "hi",
		SummonerLevel: 123,
	}

	id, err := db.InsertSummonerRecord(ctx, arg)
	assert.NoError(t, err)

	record, err := db.SelectSummonerRecordById(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, record.Name, arg.Name)
}
