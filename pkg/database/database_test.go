package database

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rank1zen/yujin/pkg/database/migrate"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	container *postgres.PostgresContainer
	dburl     string
)

const (
	dbname = "testing"
	dbuser = "yuijn"
	dbpass = "secret"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	opts := []testcontainers.ContainerCustomizer{
		postgres.WithDatabase(dbname),
		postgres.WithUsername(dbuser),
		postgres.WithPassword(dbpass),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	}

	var err error
	container, err = postgres.Run(ctx, "docker.io/postgres:13-alpine", opts...)
	if err != nil {
		log.Fatalf("running postgres container: %s", err)
	}

	dburl, err = container.ConnectionString(ctx)
	if err != nil {
		log.Fatal(err)
	}

	pgxconn, err := pgx.Connect(ctx, dburl)
	if err != nil {
		log.Fatalf("migration conn: %s", err)
	}

	err = migrate.Migrate(ctx, pgxconn)
	if err != nil {
		log.Fatalf("migrating: %s", err)
	}

	err = container.Snapshot(ctx, postgres.WithSnapshotName("test-snapshot"))
	if err != nil {
		log.Fatalf("creating snapshot: %s", err)
	}

	pgxconn.Close(ctx)

	code := m.Run()

	err = container.Terminate(ctx)
	if err != nil {
		log.Fatalf("terminating: %s", err)
	}

	os.Exit(code)
}

func setupDB(tb testing.TB) *DB {
	ctx := context.Background()

	db, err := NewDB(ctx, dburl)
	if err != nil {
		tb.Fatal(err)
	}

	tb.Cleanup(func() {
		db.Close()

		err := container.Restore(ctx)
		if err != nil {
			tb.Fatal(err)
		}
	})

	return db
}

func TestUpdateProfile(t *testing.T) {
	ctx := context.Background()

	db := setupDB(t)

	err := db.UpdateProfile(ctx, "orrange-na1")
	assert.NoError(t, err)
}
