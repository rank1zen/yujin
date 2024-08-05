package static_test

import (
	"log"
	"testing"

	"github.com/rank1zen/yujin/pkg/server/ui/static"
	"github.com/stretchr/testify/assert"
)

func TestStaticFiles(t *testing.T) {
	fs, err := static.StylesheetFiles.ReadFile("css/styles.css")
	assert.NoError(t, err)
	log.Print(string(fs))
}
