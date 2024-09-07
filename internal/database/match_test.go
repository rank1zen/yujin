package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsureMatchlist(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	db := setupDB(t)

	err := ensureMatchList(ctx, db.pool, db.riot, "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q", 0, 1)
	assert.NoError(t, err)
}
