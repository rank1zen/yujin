package database

import (
	"context"
	"log"
	"testing"

	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/riot/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfile(t *testing.T) {
	ctx := context.Background()

	db := setupDB(t)

	exists, err := db.ProfileExists(ctx, "xpzpxnzLQX12ACv3iHZfqgdA8RGZQBLCiqJVa1rfVO8Z3KRiYD7YikD2RZC5mot0YhJNKn1UDxu-Ng")
	require.NoError(t, err)
	require.True(t, exists)

	m, err := db.ProfileGetHeader(ctx, "xpzpxnzLQX12ACv3iHZfqgdA8RGZQBLCiqJVa1rfVO8Z3KRiYD7YikD2RZC5mot0YhJNKn1UDxu-Ng")
	assert.NoError(t, err)

	log.Print(m)
}

func TestProfileGetMatchListTestdata(t *testing.T) {
	ctx := context.Background()

	db := setupDB(t)

	err := matchInsert(ctx, db.pool, testdata.GetMatch("NA1_5011055088"))
	require.NoError(t, err)

	m, err := db.ProfileGetMatchList(ctx, "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q", 0, false)
	require.NoError(t, err)
	if assert.Len(t, m.List, 1) {
		log.Println(m.List[0])
	}
}

func TestProfileGetMatchListLive(t *testing.T) {
	ctx := context.Background()

	db := setupDB(t)

	_, _ = db.ProfileGetMatchList(ctx, "doublelift-na1", 0, false)
}

func TestProfileGetChampionStatList(t *testing.T) {
	ctx := context.Background()

	db := setupDB(t)

	err := matchInsert(ctx, db.pool, testdata.GetMatch("NA1_5011055088"))
	require.NoError(t, err)

	stats, err := db.ProfileGetChampionStatList(ctx, "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q", internal.SeasonAll)
	require.NoError(t, err)
	log.Print(stats)
}
