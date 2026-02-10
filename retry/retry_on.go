package retry

import (
	"context"
	"errors"
	"net"
)

// DefaultRetryOn is conservative: retries transient network-ish errors.
// It does NOT retry context cancellation/deadline by default.
func DefaultRetryOn(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, context.Canceled) {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	var ne net.Error
	if errors.As(err, &ne) {
		return ne.Timeout()
	}

	return false
}

// RetryTemporary is a helper users can compose with their own logic.
func RetryTemporary(err error) bool {
	var ne net.Error
	return errors.As(err, &ne) && (ne.Timeout())
}
