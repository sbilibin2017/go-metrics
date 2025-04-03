package context

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewContext_SignalReceived(t *testing.T) {
	ctx, cancel := NewContext()
	defer cancel()
	errCh := make(chan error, 1)
	go func() {
		time.Sleep(100 * time.Millisecond)
		err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		if err != nil {
			errCh <- err
			return
		}
		errCh <- nil
	}()
	select {
	case <-ctx.Done():
		assert.Equal(t, ctx.Err(), context.Canceled)
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Failed to send signal: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out waiting for context cancellation")
	}
}

func TestNewContext_NoSignal(t *testing.T) {
	ctx, cancel := NewContext()
	defer cancel()
	select {
	case <-ctx.Done():
		t.Fatal("Context was cancelled unexpectedly")
	case <-time.After(2 * time.Second):
	}
}
