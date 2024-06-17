package database

import (
	"context"
	"os"
	"testing"

	"github.com/rank1zen/yujin/pkg/logging"
)

var (
	testDatabaseInstance *DBTest
	riot                 RiotClient
)

func testingContext(tb testing.TB) context.Context {
	ctx := context.Background()
	return logging.WithContext(ctx, logging.NewTestLogger(tb))
}

func TestMain(m *testing.M) {
	testDatabaseInstance = MustTestInstance()
	defer testDatabaseInstance.MustClose()

	apiKey := os.Getenv("RIOT_API_KEY")
	if apiKey == "" {
		os.Exit(1)
	}

	riot = NewGolioClient(context.Background(), apiKey)

	code := m.Run()
	os.Exit(code)
}
