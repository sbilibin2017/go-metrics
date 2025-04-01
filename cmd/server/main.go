package main

import (
	"context"
	"go-metrics/cmd/server/app"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := app.ParseFlags()
	container := app.NewContainer(config)
	worker := app.NewWorker(config, container)
	server := app.NewServer(config, container, worker)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	err := server.Start(ctx)
	if err != nil {
		os.Exit(1)
	}

}
