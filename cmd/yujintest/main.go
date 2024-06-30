package main

import (
	"context"
	"net/http"
	"os"

	"github.com/rank1zen/yujin/pkg/database"
	"github.com/rank1zen/yujin/pkg/logging"
	"github.com/rank1zen/yujin/pkg/server/ui"
)

func main() {
	ctx := context.Background()

	logger := logging.NewLogger()

	db, err := database.NewDB(ctx, os.Getenv("YUJIN_TEST_DATABASE"))
	if err != nil {
		panic(err)
	}

	defer db.Close()

	ui := ui.NewUI(db, logger)
	r := ui.Routes()
	http.ListenAndServe(":8080", r)
}
