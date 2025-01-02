package main

import (
	"context"

	"github.com/jcserv/rivalslfg/internal"
	"github.com/jcserv/rivalslfg/internal/utils/log"
	"go.uber.org/zap"
)

func main() {
	logger := log.GetLogger(context.Background())
	defer logger.Sync()

	service, err := internal.NewService()
	if err != nil {
		logger.Fatal("could not create service", zap.Error(err))
	}

	if err := service.Run(); err != nil {
		logger.Fatal("could not start service", zap.Error(err))
	}
}
