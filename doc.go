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
