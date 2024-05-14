package server

import (
	"net/http"

	"github.com/rank1zen/yujin/pkg/database"
)

func (s *profilesHandler) getSummoner() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func (s *profilesHandler) fetchSummoner(riot database.RiotClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
