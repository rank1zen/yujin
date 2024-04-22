package database

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testInstance *TestInstance

func TestMain(m *testing.M) {
        testInstance = NewTestInstance()
        if testInstance.SkipDB {
                log.Printf("Skipping database tests: (%s)", testInstance.SkipReason)
        } else {
                defer testInstance.TestDB.MustClose()
        }

	code := m.Run()
	os.Exit(code)
}

func TestDatabaseConnection(t *testing.T) {
	ctx := context.Background()

        _, err := NewFromEnv(ctx, NewConfig(testInstance.TestDB.Url.String()))
        assert.NoError(t, err)
}
