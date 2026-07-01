package retry

import (
	"errors"
	"io"
	"net"
)

func RetryHTTPStatus(code int) bool {
	switch code {
	case 408, 425, 429, 500, 502, 503, 504:
		return true
	default:
		return false
	}
}

func RetryHTTPError(err error) bool {
	if err == nil {
		return false
	}

	// network errors
	var ne net.Error
	if errors.As(err, &ne) {
		return true
	}

	// EOF = broken connection
	if errors.Is(err, io.EOF) {
		return true
	}

	return false
}
