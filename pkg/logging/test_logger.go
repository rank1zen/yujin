package logging

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func NewTestLogger(tb testing.TB) *zap.Logger {
	return zaptest.NewLogger(tb, zaptest.Level(zap.DebugLevel))
}
