package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/KnutZuidema/golio"
	"github.com/rank1zen/yujin/pkg/database"
	"github.com/rank1zen/yujin/pkg/logging"
	"github.com/rank1zen/yujin/pkg/server"
)

type App struct {
	db database.DB
	gc *golio.Client
}

func (a *App) SetDatabase(db database.DB) { a.db = db }
func (a *App) GetDatabase() database.DB   { return a.db }

func (a *App) SetGolioClient(gc *golio.Client) { a.gc = gc }
func (a *App) GetGolioClient() *golio.Client   { return a.gc }

func (a *App) Run(ctx context.Context) error {
        log := logging.FromContext(ctx).Sugar()
        log.Infof("starting the server...")

	srv, err := server.NewServer("8080")
	if err != nil {
                return fmt.Errorf("failed to create server: %w", err)
	}

        r := http.NewServeMux()
	ha, err := server.NewHandler(ctx, r, a)
	if err != nil {
                return fmt.Errorf("failed to create handler: %w", err)
	}

	err = srv.ServeHTTPHandler(ctx, ha)
        if err != nil {
                return fmt.Errorf("failed to serve: %w", err)
        }

        return nil
}

func CmdRun() {
        ctx := context.Background()
        log := logging.NewLogger().Sugar()

        a := App{}
        err := database.WithNewDB(ctx, "CONN", &a)
        if err != nil {
                log.Fatalf("failed to start the flipping DB %w", err)
        }

        database.WithNewGolioClient("API KEY", &a)

        err = a.Run(ctx)
        if err != nil {
                log.Fatalf("failed to run the flipping server: %w", err)
        }
}
