package retryx_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/istvzsig/retryx"
	"github.com/istvzsig/retryx/breaker"
	"github.com/istvzsig/retryx/retry"
)

func TestWrapper_ComposesBreakerAndRetry(t *testing.T) {
	br := breaker.New(breaker.BreakerConfig{
		FailureThreshold: 1,
		OpenTimeout:      50 * time.Millisecond,
	})

	p := retry.DefaultPolicy()
	p.MaxAttempts = 2
	p.BaseDelay = 1 * time.Millisecond
	p.MaxDelay = 1 * time.Millisecond
	p.Jitter = 0
	p.RetryOn = func(error) bool { return true }

	w := retryx.Wrapper{Retry: p, Breaker: br}

	fail := errors.New("fail")
	err := w.Do(context.Background(), func(ctx context.Context) error { return fail })
	if err == nil {
		t.Fatalf("expected error")
	}

	// breaker should now be open, so wrapper should fail fast (ErrOpen)
	err2 := w.Do(context.Background(), func(ctx context.Context) error { return nil })
	if !errors.Is(err2, breaker.ErrOpen) {
		t.Fatalf("expected ErrOpen, got %v", err2)
	}
}
