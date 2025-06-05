package main

import (
	"context"
	"fmt"
	"github.com/ruslanDantsov/gophermart/internal/app"
	"github.com/ruslanDantsov/gophermart/internal/config"
	"github.com/ruslanDantsov/gophermart/internal/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	serverConfig, err := config.NewConfig(os.Args[1:])

	if err != nil {
		logger.Log.Fatal("Config initialized failed: %v", zap.Error(err))
	}

	if err := logger.Initialized(serverConfig.LogLevel); err != nil {
		logger.Log.Fatal("Logger initialized failed: %v", zap.Error(err))
	}
	defer logger.Log.Sync()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := app.NewGophermartApp(ctx, serverConfig, logger.Log)
	if err != nil {
		logger.Log.Fatal("Unable to config Server", zap.Error(err))
	}

	logger.Log.Info(fmt.Sprintf("Starting Gophermart app on %s ...", serverConfig.Address))

	if err := app.Run(ctx); err != nil {
		logger.Log.Fatal("Gophermart start failed: %v", zap.Error(err))
	}

}
