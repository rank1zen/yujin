package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	AndrewMatchID = "NA1_5011055088"
	AndrewPUUID   = "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q"
)

func TestContextTimeout(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(testingContext(t), 100*time.Millisecond)
	defer cancel()

	db := setupDB(t)

	_, err := db.getMatchPlayer(ctx, AndrewPUUID, AndrewMatchID)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestMatch_GetMatchPlayer(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(testingContext(t), 60*time.Second)
	defer cancel()

	db := setupDB(t)

	puuid := "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q"
	id := "NA1_5011055088"
	first, err := db.getMatchPlayer(ctx, puuid, id)
	require.NoError(t, err)
	require.Equal(t, time.Date(2024, time.June, 2, 0, 45, 29, 0, time.UTC), first.GameDate)

	second, err := db.getMatchPlayer(ctx, puuid, id)
	if assert.NoError(t, err) {
		assert.Equal(t, first, second)
	}
}

func TestMatch_GetMatchPlayerHasFields(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(testingContext(t), 60*time.Second)
	defer cancel()

	db := setupDB(t)

	puuid := "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q"
	id := "NA1_5011055088"
	got, err := db.getMatchPlayer(ctx, puuid, id)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, got.Assists)
	}
}

func TestMatch_GetMatchPlayerList(t *testing.T) {
	t.Parallel()

	_, cancel := context.WithTimeout(testingContext(t), 60*time.Second)
	defer cancel()

	_ = setupDB(t)

	for range []struct {
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
	}
}

// func TestGetMatchHistory(t *testing.T) {
// 	t.Parallel()
//
// 	ctx, cancel := context.WithTimeout(testingContext(t), 60*time.Second)
// 	defer cancel()
//
// 	service, db := setup2(t)
//
// 	matchIDs := []string{"NA1_5011055088", "NA1_5011024569", "NA1_5011007618", "NA1_5003179356", "NA1_5003154869"}
// 	_, err := service.insertMatches(ctx, db, matchIDs)
// 	require.NoError(t, err)
//
// 	for _, test := range []struct {
// 		puuid string
// 		start int
// 		count int
// 		want  []string
// 	}{
// 		{
// 			"0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q",
// 			0,
// 			5,
// 			matchIDs,
// 		},
// 		{
// 			"0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q",
// 			1,
// 			20,
// 			matchIDs[1:],
// 		},
// 		{
// 			"0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q",
// 			0,
// 			5,
// 			matchIDs[:5],
// 		},
// 	} {
// 		matches, err := service.getPlayerMatchHstory(ctx, db, test.puuid, test.start, test.count)
// 		got := make([]string, 0)
// 		for _, m := range matches {
// 			got = append(got, m.MatchId)
// 		}
//
// 		if assert.NoError(t, err) {
// 			assert.Equal(t, test.want, got)
// 		}
// 	}
// }
