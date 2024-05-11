package server

import (
	"net/http"

	"github.com/KnutZuidema/golio"
)

type middleware func(http.HandlerFunc, ...any) http.HandlerFunc

func chain(f http.HandlerFunc, m ...middleware) http.HandlerFunc {
	if len(m) == 0 {
		return f
	}

	return m[0](chain(f, m[1:cap(m)]...))
}

func GolioMiddleware(next http.HandlerFunc, gc *golio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
}
