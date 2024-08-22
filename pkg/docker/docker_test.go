package docker_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/pkg/docker"
	"github.com/stretchr/testify/assert"
)

var conn *pgx.Conn

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestPostgresConnect(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	var purge func()
	conn, purge = docker.NewPostgres()
	defer purge()

	conn, err := pgx.Connect(ctx, conn.Config().ConnString())
	if assert.NoError(t, err) {
		_, err := conn.Exec(ctx, "CREATE TABLE Tester (id INT)")
		assert.NoError(t, err)

		// _, err = conn.Exec(ctx, "CREATE TABLE Tester (id INT)")
		// if assert.NoError(t, err) {
		// }

		err = conn.Close(ctx)
		assert.NoError(t, err)
	}
}
