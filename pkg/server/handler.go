package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rank1zen/yujin/pkg/database"
)

type profilesHandler struct {
	db database.DB
}

type Env interface {
	GetDatabase() database.DB
	GetGolioClient() database.RiotClient
}

func NewHandler(ctx context.Context, router *http.ServeMux, env Env) (*http.ServeMux, error) {
	db := env.GetDatabase()
	if db == nil {
		return nil, fmt.Errorf("no database found in server env")
	}

	h := profilesHandler{db: db}

	gc := env.GetGolioClient()

	router.HandleFunc("GET /records/summoner/{puuid}", h.getSummoner())
	router.HandleFunc("POST /records/summoner/{puuid}/publish", h.fetchSummoner(gc))

	return router, nil
}
