package retry

import (
	"errors"
	"io"
	"net"
)

func RetryHTTPStatus(code int) bool {
	return code >= 500 || code == 429
}

// RetryHTTPError decides whether a transport-level error is retryable.
//
// This is ONLY for:
// - network failures
// - broken connections
// - transient IO issues
func RetryHTTPError(err error) bool {
	if err == nil {
		return false
	}

	// timeout / network errors
	var ne net.Error
	if errors.As(err, &ne) {
		return true
	}

	// broken connection
	if errors.Is(err, io.EOF) {
		return true
	}

	return false
}
