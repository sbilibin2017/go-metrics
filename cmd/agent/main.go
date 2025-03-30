package main

import (
	"context"
	"go-metrics/cmd/agent/app"
	"go-metrics/internal/configs"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := &configs.AgentConfig{
		Address:        ":8080",
		PollInterval:   2,
		ReportInterval: 10,
	}
	agent := app.NewMetricAgent(config)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	err := agent.Start(ctx)
	if err != nil {
		os.Exit(1)
	}
}
