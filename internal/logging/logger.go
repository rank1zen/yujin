package logging

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

const loggerKey = contextKey("logger")

func NewLogger() *zap.SugaredLogger {
	return zap.Must(zap.NewDevelopment()).Sugar()
}

func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, logger, logger)
}

func FromContext(ctx context.Context) *zap.SugaredLogger {
	logger, ok := ctx.Value(loggerKey).(*zap.SugaredLogger)
	if ok {
		return logger
	}

	return NewLogger()
}
