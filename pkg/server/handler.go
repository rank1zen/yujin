package server

import (
	"context"
	"net/http"

	"github.com/KnutZuidema/golio"
	"github.com/rank1zen/yujin/pkg/database"
)

type profilesHandler struct {
	db database.DB
}

type Env interface {
	GetDatabase() database.DB
	GetGolioClient() *golio.Client
}

func NewHandler(ctx context.Context, r *http.ServeMux, env Env) (*http.ServeMux, error) {
	h := profilesHandler{db: env.GetDatabase()}

	gc := env.GetGolioClient()

	r.HandleFunc("GET /profile/{puuid}", h.getSummoner())
	r.HandleFunc("POST /profile/{puuid}/publish", h.fetchSummoner(gc))

	return r, nil
}
