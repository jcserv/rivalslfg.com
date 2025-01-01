package log

import (
	"context"

	"github.com/jcserv/rivalslfg/internal/utils/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
)

func Init(isProd bool) *zap.Logger {
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
	logger = l.WithOptions(zap.AddCallerSkip(1))
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
