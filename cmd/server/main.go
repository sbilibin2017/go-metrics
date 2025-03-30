package main

import (
	"context"
	"go-metrics/cmd/server/app"
	"go-metrics/internal/configs"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := &configs.ServerConfig{
		Address: ":8080",
	}
	srv := app.NewServer(config)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	err := srv.Start(ctx)
	if err != nil {
		os.Exit(1)
	}

}
