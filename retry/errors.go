package retry

import (
	"errors"
	"fmt"
	"time"
)

type Error struct {
	Attempts int
	Last     error
	Elapsed  time.Duration
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("retry: failed after %d attempt(s): %v", e.Attempts, e.Last)
}

func (e *Error) Unwrap() error { return e.Last }

// Is lets errors.Is(err, target) work against the underlying error.
func (e *Error) Is(target error) bool {
	return errors.Is(e.Last, target)
}
