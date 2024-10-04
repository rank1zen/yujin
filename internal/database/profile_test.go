package database

import (
	"context"
	"log"
	"testing"

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

	m, err := db.ProfileGetMatchList(ctx, "orrange-na1", 0, false)
	require.NoError(t, err)
	if assert.Len(t, m, 1) {
		res := m[0]
		log.Println(res.ChampionIcon)
	}

}

func TestProfileGetMatchList(t *testing.T) {
	ctx := context.Background()

	db := setupDB(t)

	matches, err := db.ProfileGetMatchList(ctx, "doublelift-na1", 0, false)
	if assert.NoError(t, err) {
		log.Println(matches[0])
		log.Println(matches[1])
		log.Println(matches[2])
		log.Println(matches[3])
		log.Println(matches[4])
	}
}
