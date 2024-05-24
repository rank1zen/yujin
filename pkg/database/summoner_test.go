package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSummonerRecord(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	db := TI.NewDatabase(t)
	riot := TI.GetGolioClient()

	summQuery := db.Summoner()

	err := summQuery.FetchAndInsert(ctx, riot, "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q")
	require.NoError(t, err)

	record, err := summQuery.GetRecent(ctx, "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q")
	if assert.NoError(t, err) {
		assert.GreaterOrEqual(t, record.SummonerLevel, int32(329))
	}
}
