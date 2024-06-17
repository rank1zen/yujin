package views

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rank1zen/yujin/pkg/database"
	"github.com/rank1zen/yujin/pkg/server"
)

func disableCacheInDevMode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

type handler struct {
	db database.DB
}

type Env interface {
	DebugMode() bool
	GetDatabase() database.DB
	GetGolioClient() database.RiotClient
}

func NewRouter(ctx context.Context, env Env) (*http.ServeMux, error) {
	db := env.GetDatabase()
	if db == nil {
		return nil, fmt.Errorf("missing database")
	}

	h := handler{db: db}

	router := http.NewServeMux()
	router.Handle("/health", server.HandleHealthz(db))
	router.HandleFunc("/", h.home())
	router.HandleFunc("/profile/{puuid}", h.profile())
	return router, nil
}
