package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertMatch(t *testing.T) {
        t.Parallel()

        if TI.SkipDB {
                t.Skipf("skipping test: %s", TI.SkipReason)
        }

        ctx := context.Background()
        db := TI.GetDatabaseResource().NewDB(t)

        records := []*MatchRecord{
                {},
        }

        count, err := insertMatchRecords(db)(ctx, records)
        if assert.NoError(t, err) {
                assert.Equal(t, int64(len(records)), count)
        }
}

func TestInsertFull(t *testing.T) {
        t.Parallel()

        if TI.SkipDB {
                t.Skipf("skipping test: %s", TI.SkipReason)
        }

        ctx := context.Background()

        db := TI.GetDatabaseResource().NewDB(t)
        insertMatches := insertFullMatchRecords(db)

        batchTime := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)

        for _, test := range []struct {
                records []*FullMatchRecord
                want int64
        }{
                {
                        records: []*FullMatchRecord{
                                {
                                        Metadata: &MatchRecord{
                                                RecordDate: batchTime,
                                        },
                                },
                                {
                                        Metadata: &MatchRecord{
                                                RecordDate: batchTime,
                                        },
                                },
                        },
                        want: 1,
                },
        } {
                count, err := insertMatches(ctx, test.records)
                if assert.NoError(t, err) {
                        assert.Equal(t, test.want, count)
                }
        }
}

func TestSelectWithIds(t *testing.T) {
        t.Parallel()

        ctx := context.Background()

        if TI.SkipDB {
                t.Skipf("skipping test: %s", TI.SkipReason)
        }

        db := TI.GetDatabaseResource().NewDB(t)

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

        _, err := insertMatchRecords(db)(ctx, records)
        require.NoError(t, err)

        getMatches := getMatchRecordsMatchingIds(db)

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

                got, remainIds, err := getMatches(ctx, test.ids)
                if assert.NoError(t, err) {
                        var foundIds []string
                        for _, rec := range got {
                                foundIds = append(foundIds, rec.MatchId)
                        }

                        assert.ElementsMatch(t, want, foundIds)
                        assert.ElementsMatch(t, test.ids, append(want, remainIds...))
                }
        }
}
