package database

import (
	"os"
	"testing"
)

var TI *TestInstance

func TestMain(m *testing.M) {
	TI = NewTestInstance()
	if !TI.SkipDB {
		defer TI.GetDatabaseResource().MustClose()
	}

	code := m.Run()
	os.Exit(code)
}
