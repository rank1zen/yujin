package database

import (
	"os"
	"testing"
)

var TI TestInstance

func TestMain(m *testing.M) {
	TI = MustTestInstance()
	defer TI.MustClose()

	code := m.Run()
	os.Exit(code)
}
