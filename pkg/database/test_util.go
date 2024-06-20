package database

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

type DBTest struct {
	mu sync.Mutex

	pool      *dockertest.Pool
	container *dockertest.Resource
	url       *url.URL
	conn      *pgx.Conn
}

func NewTestInstance() (*DBTest, error) {
	ctx := context.Background()

	log.Printf("Connecting to Docker")
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Docker: %w", err)
	}

	log.Printf("Running Docker Container: %s:%s", postgresImage, postgresTag)
	container, err := pool.RunWithOptions(&dockertest.RunOptions{
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
			c.RestartPolicy = docker.RestartPolicy{Name: "no"}
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	err = container.Expire(120)
	if err != nil {
		return nil, fmt.Errorf("failed to expire database: %w", err)
	}

	connUrl := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(databaseUser, databasePassword),
		Host:     container.GetHostPort("5432/tcp"),
		Path:     databaseName,
		RawQuery: "sslmode=disable",
	}

	time.Sleep(5 * time.Second)

	conn, err := pgx.Connect(ctx, connUrl.String())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to databse: %w", err)
	}

	err = dbMigrate(ctx, conn)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate databse: %w", err)
	}

	return &DBTest{
		pool:      pool,
		container: container,
		conn:      conn,
		url:       connUrl,
	}, nil
}

func MustTestInstance() *DBTest {
	db, err := NewTestInstance()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	return db
}

func (t *DBTest) Close() error {
	err := t.conn.Close(context.Background())
	if err != nil {
		return err
	}

	err = t.pool.Purge(t.container)
	if err != nil {
		return err
	}

	return nil
}

func (t *DBTest) MustClose() {
	err := t.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

// Clones a new databse from the template, test will fatal on error
func (t *DBTest) NewDatabase(tb testing.TB) *DB {
	tb.Helper()

	name, err := t.cloneDatabase()
	if err != nil {
		tb.Fatal(err)
	}

	url := t.url.ResolveReference(&url.URL{Path: name})
	url.RawQuery = "sslmode=disable"

	ctx := context.Background()
	db, err := NewDB(ctx, url.String())
	if err != nil {
		tb.Fatalf("failed to create db :%s", err)
	}

	tb.Cleanup(func() {
		db.Close()

		ctx := context.Background()
		t.mu.Lock()
		defer t.mu.Unlock()

		_, err := t.conn.Exec(ctx, fmt.Sprintf(`DROP DATABASE IF EXISTS "%s" WITH (FORCE);`, name))
		if err != nil {
			tb.Errorf("failed to drop database %q: %s", name, err)
		}
	})

	return db
}

func (t *DBTest) NewConn(tb testing.TB) *pgx.Conn {
	tb.Helper()

	ctx := context.Background()

	name, err := t.cloneDatabase()
	if err != nil {
		tb.Fatal(err)
	}

	url := t.url.ResolveReference(&url.URL{Path: name})
	url.RawQuery = "sslmode=disable"

	conn, err := pgx.Connect(ctx, url.String())
	if err != nil {
		tb.Fatalf("failed to connect db: %s", err)
	}

	tb.Cleanup(func() {
		ctx := context.Background()

		conn.Close(ctx)

		t.mu.Lock()
		defer t.mu.Unlock()

		_, err := t.conn.Exec(ctx, fmt.Sprintf(`DROP DATABASE IF EXISTS "%s" WITH (FORCE);`, name))
		if err != nil {
			tb.Errorf("failed to drop database %q: %s", name, err)
		}
	})

	return conn
}

func (t *DBTest) NewPool(tb testing.TB) *pgxpool.Pool {
	tb.Helper()

	ctx := context.Background()

	name, err := t.cloneDatabase()
	if err != nil {
		tb.Fatal(err)
	}

	url := t.url.ResolveReference(&url.URL{Path: name})
	url.RawQuery = "sslmode=disable"

	pool, err := pgxpool.New(ctx, url.String())
	if err != nil {
		tb.Fatalf("failed to connect db: %s", err)
	}

	tb.Cleanup(func() {
		ctx := context.Background()

		pool.Close()

		t.mu.Lock()
		defer t.mu.Unlock()

		_, err := t.conn.Exec(ctx, fmt.Sprintf(`DROP DATABASE IF EXISTS "%s" WITH (FORCE);`, name))
		if err != nil {
			tb.Errorf("failed to drop database %q: %s", name, err)
		}
	})

	return pool
}

// returns name of a fresh cloned migrated db
func (t *DBTest) cloneDatabase() (string, error) {
	name, err := randomDatabaseName()
	if err != nil {
		return "", fmt.Errorf("failed to create new database name: %w", err)
	}

	// postgres does not allow parallel database creation
	t.mu.Lock()
	defer t.mu.Unlock()

	_, err = t.conn.Exec(context.Background(), fmt.Sprintf(`CREATE DATABASE "%s" WITH TEMPLATE "%s";`, name, databaseName))
	if err != nil {
		return "", fmt.Errorf("failed to clone template database: %w", err)
	}

	return name, nil
}

func randomDatabaseName() (string, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
