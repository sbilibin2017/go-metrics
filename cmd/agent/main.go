package main

import (
	"context"
	"go-metrics/cmd/agent/app"
	"go-metrics/internal/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger.Init()
	config := app.ParseFlags()
	agent := app.NewMetricAgent(config)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	err := agent.Start(ctx)
	if err != nil {
		os.Exit(1)
	}
}
