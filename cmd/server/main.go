package main

import (
	"context"
	"go-metrics/cmd/server/app"
	"go-metrics/internal/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger.Init()
	config := app.ParseFlags()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	container, err := app.NewContainer(config)
	if err != nil {
		os.Exit(1)
	}
	worker := app.NewWorker(config, container)
	server := app.NewServer(config, container, worker)
	server.Start(ctx)
	if err != nil {
		os.Exit(1)
	}

}
