package app

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/rank1zen/yujin/pkg/database"
	"github.com/rank1zen/yujin/pkg/logging"
	"github.com/rank1zen/yujin/pkg/server"
)

type App struct {
	db database.DB
	rc database.RiotClient
}

func (a *App) SetDatabase(db database.DB) { a.db = db }
func (a *App) GetDatabase() database.DB   { return a.db }

func (a *App) SetRiotClient(rc database.RiotClient) { a.rc = rc }
func (a *App) GetGolioClient() database.RiotClient   { return a.rc }

func (a *App) RunServer(ctx context.Context, port string) error {
	log := logging.FromContext(ctx).Sugar()

	srv, err := server.NewServer(port)
	if err != nil {
		return fmt.Errorf("failed to create http server: %w", err)
	}

	router := http.NewServeMux()
	// sub route
	handler, err := server.NewHandler(ctx, router, a)
	if err != nil {
		return fmt.Errorf("failed to create http handler: %w", err)
	}

	log.Infof("serving on port: %s", port)
	err = srv.ServeHTTPHandler(ctx, handler)
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

	log.Infof("starting yujin...")

	ctx = logging.WithContext(ctx, log.Desugar())

	defer func() {
		done()
		r := recover()
		if r != nil {
			log.Fatalw("application panic", "panic", r)
		}
	}()

	var cfg Config
	err := NewConfig(ctx, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	a := App{}
	defer a.Close(ctx)

	err = database.WithNewDB(ctx, &a, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to start the flipping DB %w", err)
	}

	database.WithNewGolioClient(ctx, &a, cfg.RiotApiKey)

	err = a.RunServer(ctx, cfg.Port)
	done()
	if err != nil {
		log.Fatalf("failed to start yujin: %v", err)
	}

	log.Infof("successful shutdown")
}
