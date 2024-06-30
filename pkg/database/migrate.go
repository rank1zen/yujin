package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
)

func DBMigrate(ctx context.Context, conn *pgx.Conn) error {
	migrator, err := migrate.NewMigrator(ctx, conn, "public.schema_version")
	if err != nil {
		return fmt.Errorf("could not create migrator: %w", err)
	}

	err = migrator.LoadMigrations(os.DirFS(dbMigrationsDir()))
	if err != nil {
		return fmt.Errorf("could not load migrations: %w", err)
	}

	err = migrator.Migrate(ctx)
	if err != nil {
		return fmt.Errorf("could not migrate: %w", err)
	}

	return nil
}

func dbMigrationsDir() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return filepath.Join(filepath.Dir(filename), "../../migrations")
}
