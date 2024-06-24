package database

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestSpans(t *testing.T) {
	ctx, cancel := context.WithTimeout(testingContext(t), 60*time.Second)
	defer cancel()

	db := testDatabaseInstance.NewDatabase(t)

	batch := &pgx.Batch{}
	batch.Queue("select 1")
	batch.Queue("select 2")

	br := db.pool.SendBatch(ctx, batch)

	var err error
	var n int32
	err = br.QueryRow().Scan(&n)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, n)

	err = br.QueryRow().Scan(&n)
	assert.NoError(t, err)
	assert.EqualValues(t, 2, n)

	err = br.Close()
	assert.NoError(t, err)
}
