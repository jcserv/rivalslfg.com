package log

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/jcserv/rivalslfg/internal/utils/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
)

func initSentry() error {
	return sentry.Init(sentry.ClientOptions{
		Dsn:              env.GetString("SENTRY_DSN", ""),
		AttachStacktrace: true,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		Environment:      env.GetString("ENVIRONMENT", "dev"),
	})
}

func sentryHook(entry zapcore.Entry) error {
	if entry.Level >= zapcore.ErrorLevel {
		hub := sentry.CurrentHub()
		event := sentry.NewEvent()
		event.Level = sentry.LevelError
		event.Message = entry.Message
		event.Timestamp = entry.Time
		hub.CaptureEvent(event)
	}
	return nil
}

func Init(isProd bool) *zap.Logger {
	if isProd {
		if err := initSentry(); err != nil {
			panic(err)
		}
		defer sentry.Flush(2 * time.Second)
	}

	config := zap.NewProductionConfig()
	if !isProd {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.TimeKey = "time"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.CallerKey = "caller"
		config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	}

	l, _ := config.Build()

	options := []zap.Option{zap.AddCallerSkip(1)}
	if isProd {
		sentryCore := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.RegisterHooks(core, sentryHook)
		})
		options = append(options, sentryCore)
	}

	logger = l.WithOptions(options...)
	return logger
}

func GetLogger(ctx context.Context) *zap.Logger {
	isProd := env.GetString("ENVIRONMENT", "dev") == "production"

	if logger == nil {
		return Init(isProd)
	}
	return logger
}

func Debug(ctx context.Context, msg string) {
	if !env.GetBool("DEBUG", false) {
		return
	}
	GetLogger(ctx).Debug(msg)
}

func Info(ctx context.Context, msg string) {
	GetLogger(ctx).Info(msg)
}

func Warn(ctx context.Context, msg string) {
	GetLogger(ctx).Warn(msg)
}

func Error(ctx context.Context, msg string) {
	GetLogger(ctx).Error(msg)
}

func Fatal(ctx context.Context, msg string) {
	GetLogger(ctx).Fatal(msg)
}
