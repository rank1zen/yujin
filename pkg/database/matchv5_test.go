package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelectWithIds(t *testing.T) {
        t.Parallel()

        ctx := context.Background()
        db := TI.NewDatabase(t)

        batchTime := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
        records := []*MatchRecord{
                {
                        RecordDate: batchTime,
                        MatchId: "NA1_4928428617",
                },
                {
                        RecordDate: batchTime,
                        MatchId: "NA1_4928393443",
                },
                {
                        RecordDate: batchTime,
                        MatchId: "NA1_4928353869",
                },
        }
        require.Equal(t, 1, 1)

        for _, test := range []struct {
                ids []string
                want []int
        }{
                {
                        ids: []string{},
                        want: []int{},
                },
                {
                        ids: []string{"NA1_4928428617", "NA1_4928353869", "NA1_4928393443"},
                        want: []int{0, 1, 2},
                },
                {
                        ids: []string{"NA1_4928393443"},
                        want: []int{1},
                },
        } {
                var want []string
                for _, i := range test.want {
                        want = append(want, records[i].MatchId)
                }

                got, err := db.MatchV5().GetRecords(ctx)
                if assert.NoError(t, err) {
                        var foundIds []string
                        for _, rec := range got {
                                foundIds = append(foundIds, rec.MatchId)
                        }

                        assert.ElementsMatch(t, want, foundIds)
                }
        }
}
