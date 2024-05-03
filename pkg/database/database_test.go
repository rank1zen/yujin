package database

import (
	"context"
	"os"
	"testing"
)

var TI *testInstance

func TestMain(m *testing.M) {
	TI = NewTestInstance()
	if !TI.SkipDB {
		defer TI.GetDatabaseResource().MustClose()
	}

	code := m.Run()
	os.Exit(code)
}

func TestFetchSummoner(t *testing.T) {
        t.Parallel()

        if TI.SkipDB {
                t.Skipf("skipping test: %s", TI.SkipReason)
        }

        ctx := context.Background()
        db := TI.GetDatabaseResource().NewDB(t)
}
