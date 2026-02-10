# retryx

**retryx** provides small, dependency-free primitives for building resilient Go clients:

- Exponential backoff retry with jitter
- Context-aware cancellation
- Optional per-attempt timeouts
- Circuit breaker (closed / open / half-open)
- Easy composition of retry + breaker

Designed for services calling **HTTP**, **gRPC**, or any remote dependency.

No frameworks, no dependencies, just Go.

---

## Installation

```bash
go get github.com/istvzsig/retryx
```

## Features

### Retry

- exponential backoff
- jitter support
- context cancellation respected
- per-try timeout
- custom retry decision logic
- retry hooks

### Circuit Breaker

- closed / open / half-open states
- configurable thresholds
- automatic recovery probing
- fail fast when open

### Wrapper

- Combine retry + breaker in one call wrapper.

## Quick Example

Retry + circuit breaker around a function call:

```bash
package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/istvzsig/retryx"
	"github.com/istvzsig/retryx/breaker"
	"github.com/istvzsig/retryx/retry"
)

func main() {
	br := breaker.New(breaker.BreakerConfig{
		FailureThreshold: 3,
		SuccessThreshold: 2,
		OpenTimeout:      5 * time.Second,
	})

	p := retry.DefaultPolicy()
	p.MaxAttempts = 3

	w := retryx.Wrapper{
		Retry:   p,
		Breaker: br,
	}

	ctx := context.Background()

	err := w.Do(ctx, func(ctx context.Context) error {
		return errors.New("service failed")
	})

	fmt.Println("result:", err)
}
```

## Retry Usage

### Basic retry

```bash
err := retry.Do(ctx, retry.DefaultPolicy(), func(ctx context.Context) error {
	return callService()
})
```

### Custom retry policy

```bash
p := retry.Policy{
	MaxAttempts: 5,
	BaseDelay:   100 * time.Millisecond,
	MaxDelay:    2 * time.Second,
	Jitter:      0.2,
}

err := retry.Do(ctx, p, fn)
```

### Retry hooks

```bash
p.OnRetry = func(info retry.AttemptInfo) {
	log.Printf("attempt=%d err=%v delay=%v",
		info.Attempt,
		info.Err,
		info.Delay,
	)
}
```

### Custom retry conditions

```bash
p.RetryOn = func(err error) bool {
	return retry.DefaultRetryOn(err) ||
		errors.Is(err, ErrServiceUnavailable)
}
```

### Circuit Breaker Usage

```bash
br := breaker.New(breaker.BreakerConfig{
	FailureThreshold: 5,
	SuccessThreshold: 2,
	OpenTimeout:      10 * time.Second,
})

err := br.Do(ctx, fn)
```

When open, calls immediately return:

```bash
breaker.ErrOpen
```

### Observe state changes

```bash
br := breaker.New(breaker.BreakerConfig{
	OnStateChange: func(from, to breaker.State) {
		log.Printf("breaker: %s -> %s", from, to)
	},
})
```

### HTTP Helpers

Retry helpers for HTTP:

```bash
if retry.RetryHTTPStatus(resp.StatusCode) {
	return retry.HTTPStatusError{StatusCode: resp.StatusCode}
}
```

Retryable status codes:

- 408
- 425
- 429
- 500
- 502
- 503
- 504

### Retry Error Inspection

When retries are exhausted, retryx returns a typed error:

```bash
if re, ok := retry.AsError(err); ok {
	fmt.Println("attempts:", re.Attempts)
	fmt.Println("elapsed:", re.Elapsed)
	fmt.Println("last error:", re.Last)
}
```

### Testing

Run tests:

```bash
go test ./...
```

Force fresh run:

```bash
go test -count=1 ./...
```

### Versioning

This project follows semantic versioning.

```bash
v0.x = API still evolving
v1.0 = stable API
```

### Design Goals

retryx aims to be:

- small
- dependency-free
- idiomatic Go
- production safe
- easy to adopt
- easy to remove if needed

No magic, No frameworks.
