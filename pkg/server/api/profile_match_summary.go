package api

import (
	"net/http"

	"github.com/rank1zen/yujin/pkg/database"
)

func profileMatchSummary(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
