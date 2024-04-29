package rest

import (
	"encoding/json"
	"net/http"

	"github.com/rank1zen/yujin/internal/database/summoner"
)

type GetSummonerQueries struct {
	Name  string `query:"name"`
	Puuid string `query:"puuid"`
}

func (s *Server) RegisterRouting(mux *http.ServeMux) {
        mux.HandleFunc("GET /summonerv4/", s.a)
}

// GET summoner
func (s *Server) a(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        w.Header().Set("Content-Type", "application/json")

        f := &summoner.SummonerRecordFilter{}
        records, err := s.database.SummonerV4.SelectSummmonerRecords(ctx, f)
        if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                json.NewEncoder(w).Encode(err)
        }

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(records)
}
