package app

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/rank1zen/yujin/pkg/database"
	"github.com/rank1zen/yujin/pkg/logging"
	"github.com/rank1zen/yujin/pkg/server"
	"github.com/rank1zen/yujin/pkg/server/views"
)

type App struct {
	debug bool
	db *database.DB
	riot database.RiotClient
}

func NewApp() *App {
	return &App{}
}

func (a *App) DebugMode() bool   { return a.debug }

func (a *App) SetDatabase(db *database.DB) { a.db = db }
func (a *App) GetDatabase() *database.DB   { return a.db }

func (a *App) SetRiotClient(rc database.RiotClient) { a.riot = rc }
func (a *App) GetGolioClient() database.RiotClient   { return a.riot }

func (a *App) RunServer(ctx context.Context, port string) error {
	srv, err := server.NewServer(port)
	if err != nil {
		return fmt.Errorf("failed to create http server: %w", err)
	}

	mux, err := views.NewRouter(ctx, a)
	if err != nil {
		return fmt.Errorf("failed to create views handler: %w", err)
	}

	// log.Infof("serving on port: %s", port)
	err = srv.ServeHTTPHandler(ctx, mux)
	if err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func (a *App) Close(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}

func CmdRun() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer done()

	log := logging.NewLogger().Sugar()
	// TODO: Add build tags and build id and such

	log.Infof("starting yujin...")

	defer func() {
		done()
		r := recover()
		if r != nil {
			log.Fatalw("application panic", "panic", r)
		}
	}()

	ctx = logging.WithContext(ctx, log.Desugar())

	var cfg Config
	err := NewConfig(ctx, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	app := NewApp()
	defer app.Close(ctx)

	err = database.WithNewDB(ctx, app, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to start the flipping DB %w", err)
	}

	database.WithNewGolioClient(ctx, app, cfg.RiotApiKey)

	err = app.RunServer(ctx, cfg.Port)
	done()

	if err != nil {
		log.Fatalf("failed to start yujin: %w", err)
	}

	log.Infof("successful shutdown")
}
