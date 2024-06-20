package logging

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TODO: there are some clever things we can do with logging and embedded interfaces

var (
	once   sync.Once
	logger *zap.Logger
)

type ctxKey struct{}

func NewLogger() *zap.Logger {
	once.Do(func() {
		stdout := zapcore.AddSync(os.Stdout)

		var level zapcore.Level
		levelEnv := os.Getenv("LOG_LEVEL")
		if levelEnv != "" {
			levelFromEnv, err := zapcore.ParseLevel(levelEnv)
			if err != nil {
				log.Println(fmt.Errorf("invalid level, defaulting to INFO: %w", err))
				level = zap.InfoLevel
			} else {
				level = levelFromEnv
			}
		} else {
			level = zap.InfoLevel
		}

		logLevel := zap.NewAtomicLevelAt(level)

		// prodCfg := zap.NewProductionEncoderConfig()

		devCfg := zap.NewDevelopmentEncoderConfig()
		devCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(devCfg)
		// fileEncoder := zapcore.NewJSONEncoder(prodCfg)

		// var gitRevision string
		buildInfo, _ := debug.ReadBuildInfo()
		// if ok {
		// 	for _, v := range buildInfo.Settings {
		// 		if v.Key == "vcs.revision" {
		// 			gitRevision = v.Value
		// 			return
		// 		}
		// 	}
		// }

		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, stdout, logLevel).
				With([]zapcore.Field{
					zap.String("go_version", buildInfo.GoVersion),
				}),
		)

		logger = zap.New(core)
	})

	return logger
}

func NewContext(ctx context.Context, fields ...zap.Field) context.Context {
	return context.WithValue(ctx, ctxKey{}, FromContext(ctx).With(fields...))
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
