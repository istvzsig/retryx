package retry

import (
	"errors"
	"io"
	"net"
)

// RetryHTTPError decides whether a transport-level error is retryable.
// Used for:
// - network failures
// - broken connections
// - transient IO issues
func RetryHTTPError(err error) bool {
	if err == nil {
		return false
	}

	var ne net.Error
	if errors.As(err, &ne) {
		return true
	}

	if errors.Is(err, io.EOF) {
		return true
	}

	return false
}
