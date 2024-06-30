package database

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/yujin/pkg/riot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func temp(tb testing.TB) (*pgxpool.Pool, *riot.Client) {
	return NewPool(tb), riot.NewClient()
}

func TestSummonerInserts(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(testingContext(t), 60*time.Second)
	defer cancel()

	db := setupDB(t)

	puuid := "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q"
	_, err := db.UpdateSummoner(ctx, puuid)
	assert.NoError(t, err)
}

func TestSummonerGetNewestRecord(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(testingContext(t), 60*time.Second)
	defer cancel()

	db := setupDB(t)

	puuid := "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q"
	_, err := db.UpdateSummoner(ctx, puuid)
	require.NoError(t, err)

	rec, err := db.newestSummonerRecord(ctx, puuid)
	if assert.NoError(t, err) {
		assert.Equal(t, "DXy8WEmu4Ln_5M9XKmwSZrr60we4TiYV5bQ9BVWqOecoGSc", rec.AccountId)
	}
}

func insertAndFetch(t *testing.T, ctx context.Context, db *DB) {
	puuid := "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q"
	record, err := db.UpdateSummoner(ctx, puuid)
	if assert.NoError(t, err) {
		assert.Equal(t, "DXy8WEmu4Ln_5M9XKmwSZrr60we4TiYV5bQ9BVWqOecoGSc", record.AccountId)
		assert.Equal(t, 1717304697000, record.RevisionDate)
	}
}
