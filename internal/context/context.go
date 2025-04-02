package context

import (
	"context"
	"os/signal"
	"syscall"
)

func NewContext() (context.Context, context.CancelFunc) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	return ctx, stop
}
