package retry

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

func Do(ctx context.Context, p Policy, fn func(context.Context) error) error {
	if fn == nil {
		return errors.New("retry: fn is nil")
	}
	if ctx == nil {
		return errors.New("retry: ctx is nil")
	}

	if p.MaxAttempts <= 0 {
		p.MaxAttempts = 1
	}
	if p.RetryOn == nil {
		p.RetryOn = DefaultRetryOn
	}

	start := time.Now()
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	var lastErr error
	attemptUsed := 0

	for attempt := 1; attempt <= p.MaxAttempts; attempt++ {
		attemptUsed = attempt

		tryCtx := ctx
		var cancel context.CancelFunc
		if p.TimeoutPerTry > 0 {
			tryCtx, cancel = context.WithTimeout(ctx, p.TimeoutPerTry)
		}

		err := fn(tryCtx)

		if cancel != nil {
			cancel()
		}

		if err == nil {
			return nil
		}
		lastErr = err

		// If no more attempts, break and wrap below.
		if attempt == p.MaxAttempts {
			break
		}

		// If context is done, stop early.
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// IMPORTANT: if policy says "don't retry", return raw error.
		if !p.RetryOn(err) {
			return err
		}

		delay := backoffDelay(p.BaseDelay, p.MaxDelay, attempt, p.Jitter, rnd)

		if p.OnRetry != nil {
			p.OnRetry(AttemptInfo{
				Attempt: attempt,
				Err:     err,
				Delay:   delay,
				Elapsed: time.Since(start),
			})
		}

		if delay > 0 {
			t := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				t.Stop()
				return ctx.Err()
			case <-t.C:
			}
		}
	}

	// Only wrap when we exhausted retries.
	return &Error{
		Attempts: attemptUsed,
		Last:     lastErr,
		Elapsed:  time.Since(start),
	}
}
