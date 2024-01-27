package postgresql_test

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/yujin/postgresql"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestPoolConnect(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	addr := postgresql.NewDockerResource(t)

	log, _ := zap.NewProduction()

	pool, err := postgresql.BackoffRetryPool(ctx, addr, log)
	assert.NoError(t, err)

	err = postgresql.CheckPool(ctx, pool)
	assert.NoError(t, err)
}
