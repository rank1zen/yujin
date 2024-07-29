package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	AndrewMatchID = "NA1_5011055088"
	AndrewPUUID   = "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q"
)

func TestEnsureMatchlist(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(testingContext(t), 60*time.Second)
	defer cancel()

	db := setupDB(t)

	err := db.ensureMatchlist(ctx, AndrewPUUID, 0, 1)
	assert.NoError(t, err)
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
