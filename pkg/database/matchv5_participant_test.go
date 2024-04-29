package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMatchlist(t *testing.T) {
	t.Parallel()

	if TI.SkipDB {
		t.Skipf("skipping test: %s", TI.SkipReason)
	}

	ctx := context.Background()
	db := TI.GetDatabaseResource().NewDB(t)

	records := []*MatchParticipantRecord{
		{
                        MatchId: "a",
                        Puuid: "hiu",
                },
                {
                        MatchId: "a",
                        Puuid: "hiu",
                },
                {
                        MatchId: "a",
                        Puuid: "hiu",
                },
                {
                        MatchId: "a",
                        Puuid: "hiu",
                },
	}

	count, err := insertMatchParticipantRecords(db)(ctx, records)
	require.NoError(t, err)
	require.Equal(t, int64(len(records)), count)

	for _, test := range []struct {
		puuid string
		want  []int
	}{} {
		var want []string
		for _, i := range test.want {
			want = append(want, records[i].MatchId)
		}

		got, err := getMatchlist(db)(ctx, "hi")
		if assert.NoError(t, err) {
			assert.ElementsMatch(t, want, got)
		}

	}
}
