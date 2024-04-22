package summoner

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/yujin/pkg/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSummonerRecord(t *testing.T) {
	t.Parallel()
        ctx := context.Background()

        db, err := database.NewFromEnv(ctx, database.NewConfig(testInstance.TestDB.Url.String()))
        require.NoError(t, err)

        q := NewSummonerV4Query(db)

        batchTime := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC).Truncate(time.Microsecond)
        records := []*SummonerRecordArg{
                {
                        RecordDate: batchTime,
                        Puuid: "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q",
                        AccountId: "DXy8WEmu4Ln_5M9XKmwSZrr60we4TiYV5bQ9BVWqOecoGSc",
                        SummonerId: "2xCyr5bJbp2BlMSWLRolf9_x0eSbWBay5Bam_9myXFXjZSw",
                        Name: "orrange",
                        ProfileIconId: 871,
                        SummonerLevel: 326,
                        RevisionDate: 1713500782000,
                },
                {
                        RecordDate: batchTime.Add(1 * time.Minute),
                        Puuid: "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q",
                        AccountId: "DXy8WEmu4Ln_5M9XKmwSZrr60we4TiYV5bQ9BVWqOecoGSc",
                        SummonerId: "2xCyr5bJbp2BlMSWLRolf9_x0eSbWBay5Bam_9myXFXjZSw",
                        Name: "orrange",
                        ProfileIconId: 871,
                        SummonerLevel: 326,
                        RevisionDate: 1713500782000,
                },
                {
                        RecordDate: batchTime.Add(24 * time.Hour),
                        Puuid: "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q",
                        AccountId: "DXy8WEmu4Ln_5M9XKmwSZrr60we4TiYV5bQ9BVWqOecoGSc",
                        SummonerId: "2xCyr5bJbp2BlMSWLRolf9_x0eSbWBay5Bam_9myXFXjZSw",
                        Name: "orrange",
                        ProfileIconId: 871,
                        SummonerLevel: 326,
                        RevisionDate: 1713500782000,
                },
        }
        q.InsertBatchSummonerRecord(ctx, records)

        for _, test := range []struct {
                f *SummonerRecordFilter
                want []int
        }{
                {
                        &SummonerRecordFilter{},
                        []int{0, 1, 2},
                },
        } {
                got, err := q.ReadSummonerRecords(ctx, test.f)
                if assert.NoError(t, err) {
                        for _, i := range test.want {
                                assert.Equal(t, records[i].RecordDate, got[i].RecordDate)
                        }
                }
        }
}
