package postgresql_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rank1zen/yujin/internal/postgresql"
	"github.com/rank1zen/yujin/internal/postgresql/db"
	"github.com/stretchr/testify/require"
)

var dbpool *pgxpool.Pool

func TestSummonerRecordInsertAndDelete(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	q := db.New(dbpool)

	args := db.InsertSummonerParams{
		RecordDate: postgresql.NewTimestamp(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
	}

	id, err := q.InsertSummoner(ctx, args)
	require.NoError(t, err)

	err = q.DeleteSummoner(ctx, id)
	require.NoError(t, err)
}

func TestInsertTimeout(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Nanosecond)
	defer cancel()

	q := db.New(dbpool)

	args := db.InsertSummonerParams{
		RecordDate: postgresql.NewTimestamp(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
	}

	_, err := q.InsertSummoner(ctx, args)
	require.Error(t, err)
}

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource := createContainer(pool)

	databaseUrl := fmt.Sprintf("postgres://postgres:yuyu@localhost:%s/postgres?sslmode=disable", resource.GetPort("5432/tcp"))

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		dbpool, err = pgxpool.New(ctx, databaseUrl)
		if err != nil {
			return err
		}
		return dbpool.Ping(ctx)
	}); err != nil {
		log.Fatalf("Could not pool to database: %s", err)
	}

	db, err := pgx.Connect(ctx, databaseUrl)
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	migrator, err := migrate.NewMigrator(ctx, db, "public.schema_version")
	if err != nil {
		log.Fatalf("Could not create migrator: %s", err)
	}

	if err = migrator.LoadMigrations(os.DirFS("../../db/migrations")); err != nil {
		log.Fatalf("Could not load migrations: %s", err)
	}

	if err = migrator.Migrate(ctx); err != nil {
		log.Fatalf("Could not migrate: %s", err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func createContainer(pool *dockertest.Pool) *dockertest.Resource {
	// pulls an image, creates a container based on it and runs it
	container, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "15-alpine3.18",
			Env: []string{
				"POSTGRES_PASSWORD=yuyu",
				"listen_addresses = '*'",
			},
		},
		func(hc *docker.HostConfig) {
			hc.AutoRemove = true
			hc.RestartPolicy = docker.RestartPolicy{Name: "no"}
		},
	)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	container.Expire(60)
	return container
}
