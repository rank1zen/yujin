package logging

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func NewTestLogger(tb testing.TB) *zap.Logger {
	return zaptest.NewLogger(tb, zaptest.Level(zap.DebugLevel))
}

func testingContext(tb testing.TB) context.Context {
	ctx := context.Background()
	return WithContext(ctx, NewTestLogger(tb))
}
