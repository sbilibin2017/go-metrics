package main

import (
	"context"
	"go-metrics/cmd/agent/app"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := app.ParseFlags()
	agent := app.NewMetricAgent(config)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	err := agent.Start(ctx)
	if err != nil {
		os.Exit(1)
	}
}
