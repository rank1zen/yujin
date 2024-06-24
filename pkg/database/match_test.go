package database

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatchInsertCancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(testingContext(t), 100*time.Millisecond)
	defer cancel()

	service, db := setup2(t)

	_, err := service.insertMatches(ctx, db, []string{"NA1_5011055088"})
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestMatchInsert(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(testingContext(t), 60*time.Second)
	defer cancel()

	service, db := setup2(t)

	for _, test := range []struct {
		matchIDs []string
		want     []string
	}{
		{
			[]string{"NA1_5011055088"},
			[]string{"NA1_5011055088"},
		},
		{
			[]string{"NA1_5011055088"},
			[]string{},
		},
	} {
		got, err := service.insertMatches(ctx, db, test.matchIDs)
		if assert.NoError(t, err) {
			assert.Equal(t, test.want, got)
		}
	}
}

func TestGetMatchHistory(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(testingContext(t), 60*time.Second)
	defer cancel()

	db := testDatabaseInstance.NewPool(t)
	repo := setup(t)

	matchIDs := []string{"NA1_5011055088", "NA1_5011024569", "NA1_5011007618", "NA1_5003179356", "NA1_5003154869"}
	_, err := repo.insertMatches(ctx, db, matchIDs)
	require.NoError(t, err)

	for _, test := range []struct {
		puuid string
		start int
		count int
		want  []string
	}{
		{
			"0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q",
			0,
			5,
			matchIDs,
		},
		{
			"0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q",
			1,
			20,
			matchIDs[1:],
		},
		{
			"0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q",
			0,
			5,
			matchIDs[:5],
		},
	} {
		matches, err := repo.getPlayerMatchHstory(ctx, db, test.puuid, test.start, test.count)
		got := make([]string, 0)
		for _, m := range matches {
			got = append(got, m.MatchId)
		}

		if assert.NoError(t, err) {
			assert.Equal(t, test.want, got)
		}
	}
}

func TestMatchRecords(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(testingContext(t), 120*time.Second)
	defer cancel()

	db := testDatabaseInstance.NewPool(t)
	repo := setup(t)

	ids := []string{"NA1_5011055088"}

	newIDs, err := repo.insertMatches(ctx, db, ids)
	require.NoError(t, err)
	require.Equal(t, ids, newIDs)

	want := MatchInfo{
		RecordId:     nil,
		RecordDate:   nil,
		MatchId:      "NA1_5011055088",
		GameDate:     time.Date(2024, 6, 2, 0, 0, 0, 0, time.UTC),
		GameDuration: 19*time.Minute + 25*time.Second,
		GamePatch:    "14.11.589.9418",
	}

	got, err := repo.getMatch(ctx, db, ids[0])
	if assert.NoError(t, err) {
		log.Println(got)
		assert.Equal(t, want, *got)
	}
}
