package postgres

import (
	"log"
	"os"
	"testing"

	"github.com/rank1zen/yujin/pkg/database"
)

var testInstance *database.TestInstance

func TestMain(m *testing.M) {
        testInstance = database.NewTestInstance()
        if testInstance.SkipDB {
                log.Printf("Skipping database tests: (%s)", testInstance.SkipReason)
        } else {
                defer testInstance.TestDB.MustClose()
        }

	code := m.Run()
	os.Exit(code)
}
