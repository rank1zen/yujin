package database

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rank1zen/yujin/internal/database/migrate"

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

	pgxconn.Close(ctx)

	err = container.Snapshot(ctx, postgres.WithSnapshotName("test-snapshot"))
	if err != nil {
		log.Fatalf("creating snapshot: %s", err)
	}

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

func TestGetProfileSummary(t *testing.T) {
	ctx := context.Background()

	db := setupDB(t)

	m, err := db.GetProfileSummary(ctx, "doublelift-na1")
	assert.NoError(t, err)

	log.Print(m)
}

func TestGetProfileMatchList(t *testing.T) {
	ctx := context.Background()

	db := setupDB(t)

	matches, err := db.GetProfileMatchList(ctx, "doublelift-na1", 0, true)
	if assert.NoError(t, err) {
		log.Println(matches[0])
		log.Println(matches[1])
		log.Println(matches[2])
		log.Println(matches[3])
		log.Println(matches[4])
	}
}
