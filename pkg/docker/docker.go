package docker

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const (
	databaseName     = "test-db-template"
	databaseUser     = "test-user"
	databasePassword = "testing123"

	postgresImage = "postgres"
	postgresTag   = "13-alpine"
)

var pool *dockertest.Pool

func init() {
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("connecting to docker: %v", err)
	}
}

func migrateDB(ctx context.Context, conn *pgx.Conn) error {
	migrator, err := migrate.NewMigrator(ctx, conn, "schema_version_non_default")
	if err != nil {
		return fmt.Errorf("could not create migrator: %w", err)
	}

	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("something went wrong with migrations")
	}

	migrationsDir := filepath.Join(filepath.Dir(filename), "../../migrations")
	err = migrator.LoadMigrations(os.DirFS(migrationsDir))
	if err != nil {
		return fmt.Errorf("could not load migrations: %w", err)
	}

	err = migrator.Migrate(ctx)
	if err != nil {
		return fmt.Errorf("could not migrate: %w", err)
	}

	return nil
}

func NewPostgres() (*pgx.Conn, func()) {
	ctx := context.Background()

	log.Printf("running docker container: %s:%s", postgresImage, postgresTag)
	container, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: postgresImage,
			Tag:        postgresTag,
			Env: []string{
				"POSTGRES_DB=" + databaseName,
				"POSTGRES_USER=" + databaseUser,
				"POSTGRES_PASSWORD=" + databasePassword,
			},
		},
		func(c *docker.HostConfig) {
			c.AutoRemove = true
			c.RestartPolicy = docker.NeverRestart()
		},
	)

	if err != nil {
		log.Fatalf("starting psql container: %v", err)
	}

	container.Expire(120) // error handling?

	connUrl := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(databaseUser, databasePassword),
		Host:     container.GetHostPort("5432/tcp"),
		Path:     databaseName,
		RawQuery: "sslmode=disable",
	}

	log.Print(connUrl.String())
	time.Sleep(10 * time.Second)

	conn, err := pgx.Connect(ctx, connUrl.String())
	if err != nil {
		log.Fatalf("pgx: %v", err)
	}

	return conn, func() {
		conn.Close(ctx)
		pool.Purge(container)
	}
}
