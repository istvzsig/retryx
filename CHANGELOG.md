# Changelog

All notable changes to `retryx` will be documented in this file.

The format is based on Keep a Changelog, and this project follows Semantic Versioning.

---

## [v1.0.0] - 2026-07-01

### 🎉 First stable release

This is the first stable version of retryx.  
The core API is now frozen and production-ready.

---

## Added

### Retry system

- Exponential backoff retry engine
- Jitter support for retry delays
- Context-aware cancellation
- Per-attempt timeout support
- Custom retry policy configuration
- Retry hooks for observability

### Circuit Breaker

- Closed / Open / Half-Open states
- Configurable failure threshold
- Configurable success threshold
- Automatic recovery after open timeout
- State transition hooks

### Composition layer

- Unified `Wrapper` for retry + circuit breaker
- Works with any function signature `(context.Context) error`
- Clean separation of retry and breaker logic

### HTTP support patterns

- Standardized retry decision patterns for HTTP calls
- Retryable network error handling helper (`RetryHTTPError`)
- Example integration with `net/http`

### Examples

- Failure simulation demo (retry + breaker behavior)
- HTTP client integration example

---

## Stability guarantees (v1.0)

From v1.0 onwards:

- ❌ No breaking changes without major version bump
- ❌ No removal of exported APIs
- ❌ No signature changes in core components
- ✔ Only backward-compatible additions allowed

---

## Core API (stable surface)

The following APIs are now stable:

### retry

- `retry.Do(...)`
- `retry.DefaultPolicy(...)`
- `retry.RetryHTTPError(...)`

### breaker

- `breaker.New(...)`
- `(*Breaker).Do(...)`
- `(*Breaker).State(...)`

### wrapper

- `retryx.Wrapper`
- `(*Wrapper).Do(...)`

---

## Design philosophy

- No frameworks
- No hidden magic
- Explicit behavior only
- Small, composable primitives
- Idiomatic Go

---

## Known limitations

- No built-in metrics (planned for v1.1)
- No OpenTelemetry integration yet
- No retry budget system yet

---

## Upgrade path

Future versions may add:

- Metrics & observability hooks
- Retry budget system
- OpenTelemetry integration
- Enhanced logging hooks

All additions will remain backward compatible.

---

## Summary

`retryx v1.0` is a stable, production-ready resilience library for Go:

- Retry logic
- Circuit breaker
- Composition wrapper
- HTTP-friendly patterns

Ready for production use.
