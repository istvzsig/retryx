# retryx

> Lightweight, dependency-free resilience primitives for Go.

retryx helps you build **production-grade resilient clients** with:

- retries (exponential backoff + jitter)
- circuit breakers
- context-aware cancellation
- optional per-attempt timeouts
- composable retry + breaker wrapper

Designed for real backend systems calling HTTP, gRPC, or any unreliable dependency.

---

## Why retryx?

In distributed systems, failures are normal:

- APIs fail intermittently
- services degrade under load
- network calls timeout
- downstream dependencies become unstable

Without proper resilience, these failures can cascade and bring down entire systems.

retryx provides simple building blocks to prevent that.

---

## Features

### Retry

- Exponential backoff
- Jitter support
- Context cancellation support
- Per-attempt timeout
- Custom retry logic
- Retry hooks

### Circuit Breaker

- Closed / Open / Half-open states
- Configurable thresholds
- Automatic recovery probing
- Fast failure when open
- State change hooks

### Composition

- Combine retry + circuit breaker in one wrapper
- Works with any function (HTTP, gRPC, DB, etc.)

---

## Installation

```bash
go get github.com/istvzsig/retryx
```

---

## Quick Example

### Retry + Circuit Breaker

```go
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

	policy := retry.DefaultPolicy()
	policy.MaxAttempts = 3

	client := retryx.Wrapper{
		Retry:   policy,
		Breaker: br,
	}

	ctx := context.Background()

	err := client.Do(ctx, func(ctx context.Context) error {
		return errors.New("service failed")
	})

	fmt.Println("result:", err)
}
```

---

## Retry Example

```go
err := retry.Do(ctx, retry.DefaultPolicy(), func(ctx context.Context) error {
	return callService()
})
```

---

## Circuit Breaker Example

```go
br := breaker.New(breaker.BreakerConfig{
	FailureThreshold: 5,
	SuccessThreshold: 2,
	OpenTimeout:      10 * time.Second,
})

err := br.Do(ctx, fn)
```

---

## HTTP Retry Helper

```go
if retry.RetryHTTPStatus(resp.StatusCode) {
	return retry.HTTPStatusError{StatusCode: resp.StatusCode}
}
```

---

## Architecture

```text
Client
  ↓
Retry (backoff + jitter)
  ↓
Circuit Breaker
  ↓
External Service (HTTP / gRPC / DB)
```

---

## Design Goals

- Small and readable
- No dependencies
- Idiomatic Go
- Easy to embed in services
- Easy to remove if needed
- Explicit behavior (no hidden magic)

---

## Error Inspection

```go
if re, ok := retry.AsError(err); ok {
	fmt.Println("attempts:", re.Attempts)
	fmt.Println("elapsed:", re.Elapsed)
	fmt.Println("last error:", re.Last)
}
```

---

## Testing

```bash
go test ./...
```

Run fresh:

```bash
go test -count=1 ./...
```

---

## Versioning

- v0.x → evolving API
- v1.0 → stable API

---

## Use cases

retryx is ideal for:

- HTTP clients
- gRPC clients
- microservices
- background workers
- distributed systems
- external API integration

---

## Summary

retryx provides simple, composable resilience primitives for real-world Go systems.

No frameworks. No dependencies. No magic.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
