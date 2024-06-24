package database

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/yujin/pkg/logging"
	"github.com/rank1zen/yujin/pkg/riot"
	"github.com/stretchr/testify/assert"
)

var testDatabaseInstance *DBTest

func testingContext(tb testing.TB) context.Context {
	ctx := context.Background()
	return logging.WithContext(ctx, logging.NewTestLogger(tb))
}

func setup(tb testing.TB) *service {
	return &service{
		riot: riot.NewClient(),
	}
}

func setup2(tb testing.TB) (*service, *pgxpool.Pool) {
	return setup(tb), testDatabaseInstance.NewPool(tb)
}

func TestMain(m *testing.M) {
	testDatabaseInstance = MustTestInstance()
	defer testDatabaseInstance.MustClose()

	code := m.Run()
	os.Exit(code)
}

func TestUpdateMatchHistory(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(testingContext(t), 60*time.Second)
	defer cancel()

	db := testDatabaseInstance.NewDatabase(t)

	err := db.UpdateMatchHistory(ctx, "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q")
	assert.NoError(t, err)
}
