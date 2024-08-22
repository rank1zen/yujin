package migrate

import (
	"context"
	"os"
	"testing"

	"github.com/rank1zen/yujin/pkg/docker"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	ret := m.Run()
	os.Exit(ret)
}

func TestMigrate(t *testing.T) {
	ctx := context.Background()

	conn, purge := docker.NewPostgres()
	defer purge()

	err := Migrate(ctx, conn)
	assert.NoError(t, err)
}
