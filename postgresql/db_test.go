package postgresql_test

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/yujin/postgresql"
	"github.com/stretchr/testify/assert"
)

func TestPoolConnect(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	addr := postgresql.NewDockerResource(t)


	pool, err := postgresql.BackoffRetryPool(ctx, addr)
	assert.NoError(t, err)

	err = postgresql.CheckPool(ctx, pool)
	assert.NoError(t, err)

	err = postgresql.Migrate(ctx, pool)
	assert.NoError(t, err)
}

func TestMigrate(t *testing.T) {
	_, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
}
