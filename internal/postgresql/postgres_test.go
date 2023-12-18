package postgresql_test

import (
	"log"
	"testing"

	"github.com/rank1zen/yujin/internal/postgresql"
)

func TestParseUUID(t *testing.T) {
	_, err := postgresql.ParseUUID("string")
	if err == nil {
		log.Fatal("Expected invalid UUID")
	}
	_, err = postgresql.ParseUUID("40e6215d-b5c6-4896-987c-f30f3678f608")
	if err != nil {
		log.Fatalf("Expected success (STRING TO UUID) but got error: %s", err)
	}
}
