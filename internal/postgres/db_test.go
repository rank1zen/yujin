package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/yujin/internal/postgres"
	"github.com/stretchr/testify/assert"
)

func TestConnectPool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	addr := NewDockerResource(t)


	_, err := postgres.NewBackoffPool(ctx, addr)
	assert.NoError(t, err)
}
