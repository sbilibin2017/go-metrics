package main

import (
	"go-metrics/cmd/server/app"
	"go-metrics/internal/context"
	"go-metrics/internal/logger"
	"os"
)

func main() {
	logger.Init()
	config := app.ParseFlags()
	container, err := app.NewContainer(config)
	if err != nil {
		os.Exit(1)
	}
	worker := app.NewWorker(config, container)
	server := app.NewServer(config, container, worker)
	ctx, stop := context.NewContext()
	defer stop()
	err = server.Start(ctx)
	if err != nil {
		os.Exit(1)
	}

}
