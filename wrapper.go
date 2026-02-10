// Package retryx provides small, dependency-free primitives for retry and circuit breaking.
//
// Typical usage wraps an outbound call with both retry and a circuit breaker:
//
//  br := breaker.New(breaker.BreakerConfig{FailureThreshold: 3})
//  w := retryx.Wrapper{Retry: retry.DefaultPolicy(), Breaker: br}
//
//  err := w.Do(ctx, func(ctx context.Context) error {
//      // call remote dependency
//      return nil
//  })

package retryx

import (
	"context"

	"github.com/istvzsig/retryx/breaker"
	"github.com/istvzsig/retryx/retry"
)

type Wrapper struct {
	Retry   retry.Policy
	Breaker *breaker.Breaker
}

func (w Wrapper) Do(ctx context.Context, fn func(context.Context) error) error {
	call := fn
	if w.Breaker != nil {
		call = func(c context.Context) error {
			return w.Breaker.Do(c, fn)
		}
	}

	p := w.Retry
	if p.MaxAttempts == 0 {
		p = retry.DefaultPolicy()
	}
	return retry.Do(ctx, p, call)
}
