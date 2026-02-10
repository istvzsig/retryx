package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDo_SucceedsAfterRetries(t *testing.T) {
	ctx := context.Background()

	var calls int
	p := DefaultPolicy()
	p.MaxAttempts = 3
	p.BaseDelay = 1 * time.Millisecond
	p.MaxDelay = 2 * time.Millisecond
	p.Jitter = 0
	p.RetryOn = func(err error) bool { return true }

	err := Do(ctx, p, func(ctx context.Context) error {
		calls++
		if calls < 3 {
			return errors.New("nope")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_StopsWhenRetryOnFalse(t *testing.T) {
	ctx := context.Background()

	var calls int
	want := errors.New("fatal")
	p := DefaultPolicy()
	p.MaxAttempts = 5
	p.RetryOn = func(err error) bool { return false }

	err := Do(ctx, p, func(ctx context.Context) error {
		calls++
		return want
	})
	if !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RespectsContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	p := DefaultPolicy()
	p.MaxAttempts = 3

	err := Do(ctx, p, func(ctx context.Context) error {
		return errors.New("x")
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
