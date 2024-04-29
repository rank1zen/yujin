package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSummonerRecord(t *testing.T) {
	t.Parallel()
        ctx := context.Background()

        db := TI.GetDatabaseResource().NewDB(t)

        batchTime := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC).Truncate(time.Microsecond)
        records := []*SummonerRecord{
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

        _, err := insertSummonerRecords(db)(ctx, records)
        require.NoError(t, err)

        getSummoners := getSummonerRecords(db, nil)
        for _, test := range []struct {
                f []RecordFilter
                want []int
        }{
                {
                        []RecordFilter{
                                {
                                        Field: "name",
                                        Value: "orrange",
                                },
                        },
                        []int{0, 1, 2},
                },
        } {
                a := test.f
                got, err := getSummoners(ctx, a...)
                if assert.NoError(t, err) {
                        for _, i := range test.want {
                                assert.Equal(t, records[i].RecordDate, got[i].RecordDate)
                        }
                }
        }
}
