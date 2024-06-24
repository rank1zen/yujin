package riot

import (
	"context"
	"testing"

	"github.com/rank1zen/yujin/pkg/logging"
)

func setup(tb testing.TB) *Client {
	return NewClient(WithToken2())
}

func testingContext(tb testing.TB) context.Context {
	ctx := context.Background()
	return logging.WithContext(ctx, logging.NewTestLogger(tb))
}
