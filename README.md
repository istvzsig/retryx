# retryx

![CI](https://github.com/istvzsig/retryx/actions/workflows/ci.yml/badge.svg)
![Go Version](https://img.shields.io/badge/go-1.22-blue)
![License](https://img.shields.io/badge/license-MIT-green)

retryx is a lightweight Go library for building **resilient distributed systems**.

It provides **retry + circuit breaker primitives without frameworks, hidden magic, or dependencies**.

Designed for real backend systems calling HTTP, gRPC, databases, or any unreliable dependency.

---

## Why retryx?

In distributed systems, failures are normal:

- APIs fail intermittently
- Services degrade under load
- Network calls timeout
- Downstream dependencies become unstable

Without proper resilience, these failures can cascade and bring down entire systems.

retryx provides simple building blocks to prevent this.

---

## Features

### Retry

- Exponential backoff
- Jitter support
- Context-aware cancellation
- Per-attempt timeout
- Custom retry logic
- Retry hooks
- HTTP retry helpers

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

## Quick Start

### Retry + Circuit Breaker

```go
package main

import (
	"context"
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
		return callService(ctx)
	})

	fmt.Println("result:", err)
}
```

---

## Retry Example

```go
err := retry.Do(ctx, retry.DefaultPolicy(), func(ctx context.Context) error {
	return callService(ctx)
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

err := br.Do(ctx, func(ctx context.Context) error {
	return callService(ctx)
})
```

---

## HTTP Retry Helper

Retry decision should be made at the call site:

```go
import "errors"

err := client.Do(ctx, func(ctx context.Context) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://example.com",
		nil,
	)
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// retryable HTTP statuses
	if resp.StatusCode >= 500 || resp.StatusCode == 429 {
		return errors.New("retryable http status")
	}

	return nil
})
```

---

## HTTP Client Example

A complete example wrapping the standard `net/http` client with retry and circuit breaker is included.

```bash
go run ./examples/httpclient
```

The example demonstrates:

- context-aware HTTP requests
- retryable HTTP status handling
- retryable network errors
- retry + circuit breaker composition

---

## Failure Simulation Demo

Simulates an unstable upstream service and demonstrates:

- retry with exponential backoff
- circuit breaker opening after repeated failures
- automatic recovery after the open timeout

```bash
go run ./examples/failure_demo
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

Race detector test:

```bash
go test -race ./...
```

---

## Versioning

- v0.x → evolving API
- v1.0 → stable API

---

- resilient HTTP clients
- gRPC clients
- microservices
- background workers
- distributed systems
- database access
- third-party API integrations

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## Summary

retryx provides simple, composable resilience primitives for real-world Go systems.

No frameworks. No dependencies. No magic.

## v1.0 API Stability

Once v1.0 is released:

- No breaking changes without major version bump
- Retry and breaker APIs will remain stable
- Only additive changes allowed (backwards compatible)
