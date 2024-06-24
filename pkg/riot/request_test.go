package riot

import (
	"log"
	"testing"
)

func TestRequest(t *testing.T) {
	req := NewRequest(
		WithToken2(),
		WithURL("hihi"),
		WithQuery("val", "2"),
	)
	log.Println(req.token, req.url, req.query)
}
