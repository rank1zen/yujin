package server

import (
	"net/http"

	"github.com/KnutZuidema/golio"
)

func (s *profilesHandler) getSummoner() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		r.ParseForm()

		_, err := s.db.SummonerV4().GetRecords(ctx)
		if err != nil {

		}

		// how tf to use templ?
	}
}

func (s *profilesHandler) fetchSummoner(gc *golio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := r.PathValue("puuid")
		if id == "" {
			w.Write([]byte("Error here."))
			return
		}

		var puuid string
		err := s.db.FetchAndInsertSummoner(ctx, gc, puuid)
		if err != nil {
			w.Write([]byte("Error here."))
			return
		}
	}
}
