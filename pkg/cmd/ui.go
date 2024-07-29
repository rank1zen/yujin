package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/rank1zen/yujin/pkg/database"
	"github.com/rank1zen/yujin/pkg/logging"
	"github.com/rank1zen/yujin/pkg/server/ui"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func UiCmd(ctx context.Context) *cobra.Command {
	var port int

	return &cobra.Command{
		Use:   "ui",
		Short: "Runs the ui web server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			port = 4001
			if os.Getenv("YUJIN_PORT") != "" {
				port, _ = strconv.Atoi(os.Getenv("YUJIN_PORT"))
			}

			logger := logging.NewLogger()
			defer func() { logger.Sync() }()

			db, err := database.NewDB(ctx, os.Getenv("YUJIN_POSTGRES_POOL_URL"))
			if err != nil {
				return err
			}

			defer db.Close()

			ui := ui.Routes(db, logger)
			go func() { http.ListenAndServe(fmt.Sprintf(":%d", port), ui) }()

			logger.Info("started ui", zap.Int("port", port))

			<-ctx.Done()

			return nil
		},
	}
}
