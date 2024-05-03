package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rank1zen/yujin/pkg/database"
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
        r.FormValue("hi")
        r.ParseForm()

        // bind some data this is kinda sus rn because we dont know how we are going to read/validate/what type of data is it even gonna be
        bindQuery(r, )

        var f []database.RecordFilter

        // call some service
        records, err := s.db.Record.SummonerV4.GetRecords(ctx, f...)
        if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                json.NewEncoder(w).Encode(err)
        }

        // we are most definitely responding in json right boys?
        data, err := json.Marshal(records)
        if err != nil {
                // Handle the error when marshalling json please
        }

        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "%s", data)
}

// we want bind the queries to some struct, these could have some funny types
// ofc we will have strings and ints, but also dates lol
// if the types dont match exactly what we want then we probably want to just ignore the values
func bindQuery(r *http.Request) {
        err := r.ParseForm()
        if err != nil {

        }

        ok := r.Form["hi"]

        for key, val := range r.Form {
        }
}
