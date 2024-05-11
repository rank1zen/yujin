package database

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var TI TestInstance

func TestMain(m *testing.M) {
	TI = MustTestInstance()
	defer TI.MustClose()

	code := m.Run()
	os.Exit(code)
}

func TestFetchSummoner(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := TI.NewDatabase(t)
	gc := TI.GetGolioClient()

	err := db.FetchAndInsertSummoner(ctx, gc, "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q")
	assert.NoError(t, err)
}
