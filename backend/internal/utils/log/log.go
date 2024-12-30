package log

import (
	"context"

	"github.com/jcserv/rivalslfg/internal/utils/env"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func Init(isProd bool) *zap.Logger {
	var logger *zap.Logger
	if isProd {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	return logger
}

func GetLogger(ctx context.Context) *zap.Logger {
	isProd := env.GetString("ENVIRONMENT", "dev") == "production"

	if logger == nil {
		return Init(isProd)
	}
	return logger
}

func Error(ctx context.Context, msg string) {
	GetLogger(ctx).Error(msg)
}

func Fatal(ctx context.Context, msg string) {
	GetLogger(ctx).Fatal(msg)
}

func Info(ctx context.Context, msg string) {
	GetLogger(ctx).Info(msg)
}
