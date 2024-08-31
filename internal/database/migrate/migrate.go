package migrate

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
)

//go:embed migrations/*.sql
var embeddedQueries embed.FS

func Migrate(ctx context.Context, conn *pgx.Conn) error {
	migrator, err := migrate.NewMigrator(ctx, conn, "public.schema_version")
	if err != nil {
		return fmt.Errorf("creating migrator: %w", err)
	}

	migrations, _ := fs.Sub(embeddedQueries, "migrations")
	err = migrator.LoadMigrations(migrations)
	if err != nil {
		return fmt.Errorf("loading migrations: %w", err)
	}

	err = migrator.Migrate(ctx)
	if err != nil {
		return fmt.Errorf("migrating: %w", err)
	}

	return nil
}
