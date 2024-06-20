package riot

import (
	"context"
	"net/http"
	"testing"

	"github.com/rank1zen/yujin/pkg/logging"
	"github.com/stretchr/testify/assert"
)

func testingContext(tb testing.TB) context.Context {
	ctx := context.Background()
	return logging.WithContext(ctx, logging.NewTestLogger(tb))
}

func TestMatchlist(t *testing.T) {
	t.Parallel()

	ctx := testingContext(t)
	c := &http.Client{}

	ids, err := listByPuuid(ctx, c, "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q", 0, 5)
	if assert.NoError(t, err) {
		assert.Len(t, ids, 5)
	}
}
