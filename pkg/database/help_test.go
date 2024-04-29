package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRef(t *testing.T) {

        batchTime := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC).Truncate(time.Microsecond)
        records := []SummonerRecord{
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

        fields, rows, err := ExtractStructSlice(records)
        if assert.NoError(t, err) {
                assert.ElementsMatch(t, []string{"record_date"}, fields)
                t.Log(fields)
                t.Log(rows)
        }

	_ = []string{"record_date", "account_id", "profile_icon_id", "revision_date", "name", "summoner_id", "puuid", "summoner_level"}
}
