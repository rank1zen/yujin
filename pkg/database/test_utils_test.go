package database

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseConnection(t *testing.T) {
        t.Parallel()

        if TI.SkipDB {
                t.Skipf("skipping test: %s", TI.SkipReason)
        }

	ctx := context.Background()
        db := TI.GetDatabaseResource().NewDB(t)

        var count int
        err := db.QueryRow(ctx, "SELECT COUNT(*) FROM MatchRecords").Scan(&count)
        if assert.NoError(t, err) {
                assert.Equal(t, 0, count)
        }

        err = db.QueryRow(ctx, "SELECT COUNT(*) FROM SummonerRecords").Scan(&count)
        if assert.NoError(t, err) {
                assert.Equal(t, 0, count)
        }

        err = db.QueryRow(ctx, "SELECT COUNT(*) FROM MatchBanRecords").Scan(&count)
        if assert.NoError(t, err) {
                assert.Equal(t, 0, count)
        }
}

func TestDatabaseClone(t *testing.T) {
        t.Parallel()

        if TI.SkipDB {
                t.Skipf("skipping test: %s", TI.SkipReason)
        }

        _ = context.Background()

        a := func(st *testing.T) {
                st.Parallel()

                ctx := context.Background()
                db := TI.GetDatabaseResource().NewDB(st)

                batchTime := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
                _, err := db.Exec(ctx, "INSERT INTO SummonerRecords (record_date) VALUES ($1)", batchTime)
                require.NoError(st, err)

                var count int
                err = db.QueryRow(ctx, "SELECT COUNT(*) FROM SummonerRecords").Scan(&count)
                if assert.NoError(t, err) {
                        assert.Equal(t, 1, count)
                }
        }

        for i := range 10 {
                r := t.Run(fmt.Sprintf("Subtest %d", i), a)
                assert.True(t, r)
        }
}
