package logging

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

// TODO: there are some clever things we can do with logging and embedded interfaces

var (
	once   sync.Once
	logger *zap.Logger
)

type ctxKey struct{}

func NewLogger() *zap.Logger {
        return zap.Must(zap.NewDevelopment())
}

func FromContext(ctx context.Context) *zap.Logger {
        if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
                return l
        } else if l := logger; l != nil {
                return l
        } else {
                return zap.NewNop()
        }
}

func WithContext(ctx context.Context, lg *zap.Logger) context.Context {
        if lp, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
                if lp == lg {
                        return ctx
                }
        }

	return context.WithValue(ctx, ctxKey{}, lg)
}
