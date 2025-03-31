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
	srv := app.NewServer(config)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	err := srv.Start(ctx)
	if err != nil {
		os.Exit(1)
	}

}
