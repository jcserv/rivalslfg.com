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
	var config zap.Config
	if isProd {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.TimeKey = "time"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.CallerKey = "caller"
		config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	}

	logger, _ := config.Build(zap.AddCallerSkip(1))
	return logger
}

func GetLogger(ctx context.Context) *zap.Logger {
	isProd := env.GetString("ENVIRONMENT", "dev") == "production"

	if logger == nil {
		return Init(isProd)
	}
	return logger
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
