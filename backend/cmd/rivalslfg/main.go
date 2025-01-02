package main

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/jcserv/rivalslfg/internal"
	"github.com/jcserv/rivalslfg/internal/utils/env"
	"github.com/jcserv/rivalslfg/internal/utils/log"
	"go.uber.org/zap"
)

func main() {
	logger := log.GetLogger(context.Background())
	defer logger.Sync()

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              env.GetString("SENTRY_DSN", ""),
		AttachStacktrace: true,
	}); err != nil {
		logger.Fatal("could not initialize sentry", zap.Error(err))
	}

	defer sentry.Flush(2 * time.Second)
	sentry.CaptureMessage("It works!")

	service, err := internal.NewService()
	if err != nil {
		logger.Fatal("could not create service", zap.Error(err))
	}

	if err := service.Run(); err != nil {
		logger.Fatal("could not start service", zap.Error(err))
	}
}
