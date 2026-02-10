package breaker

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestBreaker_OpensAfterFailures(t *testing.T) {
	b := New(BreakerConfig{
		FailureThreshold: 2,
		OpenTimeout:      50 * time.Millisecond,
	})

	fail := errors.New("fail")

	_ = b.Do(context.Background(), func(ctx context.Context) error { return fail })
	_ = b.Do(context.Background(), func(ctx context.Context) error { return fail })

	if st := b.State(); st != Open {
		t.Fatalf("expected Open, got %v", st)
	}

	err := b.Do(context.Background(), func(ctx context.Context) error { return nil })
	if !errors.Is(err, ErrOpen) {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestBreaker_HalfOpenThenClosesOnSuccess(t *testing.T) {
	b := New(BreakerConfig{
		FailureThreshold: 1,
		SuccessThreshold: 2,
		OpenTimeout:      10 * time.Millisecond,
	})

	fail := errors.New("fail")
	_ = b.Do(context.Background(), func(ctx context.Context) error { return fail })

	time.Sleep(12 * time.Millisecond)

	_ = b.Do(context.Background(), func(ctx context.Context) error { return nil })
	if st := b.State(); st != HalfOpen {
		t.Fatalf("expected HalfOpen, got %v", st)
	}

	_ = b.Do(context.Background(), func(ctx context.Context) error { return nil })
	if st := b.State(); st != Closed {
		t.Fatalf("expected Closed, got %v", st)
	}
}
