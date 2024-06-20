package database

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/yujin/pkg/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchMatchlist(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(testingContext(t), 60*time.Second)
	defer cancel()

	db := testDatabaseInstance.NewPool(t)

	puuid := "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q"
	req := NewMatchRequest{Puuid: puuid, Start: 0, Count: 4}

	newIDs, err := fetchNewMatches(ctx, db, riot, req)
	require.NoError(t, err)
	require.Len(t, newIDs, req.Count)

	first, err := updateMatchHistory(ctx, db, riot, puuid)
	require.NoError(t, err)
	require.Len(t, first, req.Count)

	second, err := updateMatchHistory(ctx, db, riot, puuid)
	if assert.NoError(t, err) {
		assert.Less(t, len(second), len(first))
	}

	// LewxuBYhgAt9KI1x8vSnmA63Kg3_hN8rabMVWzK06Zg6T-j1pz-gp4qtY4jOtMzNix90I3D2GylSCA
}

func TestMatchHistory(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(testingContext(t), 120*time.Second)
	defer cancel()

	db := testDatabaseInstance.NewPool(t)

	matchIDs := [][]string{
		{
			"NA1_5011055088", "NA1_5011024569", "NA1_5011007618", "NA1_5003179356", "NA1_5003154869",
			"NA1_5000922043", "NA1_4997891656", "NA1_4997274155", "NA1_4995777967", "NA1_4995159232",
			"NA1_4994471478", "NA1_4993718318", "NA1_4993177950", "NA1_4993150709", "NA1_4993092016",
			"NA1_4991698080", "NA1_4991095671", "NA1_4991079469", "NA1_4991062852", "NA1_4991033673",
		},
		{
			"NA1_5021422737", "NA1_5021388717", "NA1_5021373818", "NA1_5020853384", "NA1_5020832606",
			"NA1_5020810810", "NA1_5020057215", "NA1_5019976769", "NA1_5019955916", "NA1_5019933370",
			"NA1_5019916474", "NA1_5019187873", "NA1_5019155528", "NA1_5019131994", "NA1_5019112396",
			"NA1_5019109886", "NA1_5015788779", "NA1_5015769535", "NA1_5015755984", "NA1_5015740407",
		},
		{
			"NA1_5021293538", "NA1_5021281021", "NA1_5021258185", "NA1_5021233370", "NA1_5020885153",
			"NA1_5020853384", "NA1_5020820862", "NA1_5020781173", "NA1_5020759046", "NA1_5018964202",
			"NA1_5018607928", "NA1_5018558505", "NA1_5016611779", "NA1_5016557236", "NA1_5016527451",
			"NA1_5016501726", "NA1_5016471341", "NA1_5016446580", "NA1_5016427496", "NA1_5016404354",
		},
	}

	for _, ids := range matchIDs {
		_, err := insertMatches(ctx, db, riot, ids)
		require.NoError(t, err)
	}

	for _, test := range []struct {
		puuid string
		start int
		count int
		want  []string
	}{
		{
			"0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q",
			0,
			20,
			matchIDs[0],
		},
		{
			"0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q",
			1,
			20,
			matchIDs[0][1:],
		},
		{
			"0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q",
			0,
			5,
			matchIDs[0][:5],
		},
	} {
		matches, err := getPlayerMatchHistory(ctx, db, test.puuid, test.start, test.count)
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

	logger := logging.FromContext(ctx).Sugar()

	db := testDatabaseInstance.NewPool(t)

	ids := []string{"NA1_5011055088"}
	puuid := "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q"

	newIDs, err := insertMatches(ctx, db, riot, ids)
	require.NoError(t, err)
	require.Equal(t, ids, newIDs)

	match, err := getPlayerMatchHistory(ctx, db, puuid, 0, 1)
	require.NoError(t, err)
	assert.Len(t, match, 1)

	logger.Info(match[0])
}
