package retry

import "time"

type Policy struct {
	// MaxAttempts includes the first try. Must be >= 1.
	MaxAttempts int

	BaseDelay time.Duration
	MaxDelay  time.Duration

	// Jitter is 0..1 (fraction of delay). 0 = no jitter.
	Jitter float64

	// TimeoutPerTry applies a timeout to each attempt (optional).
	TimeoutPerTry time.Duration

	// RetryOn decides whether an error should be retried.
	// If nil, DefaultRetryOn is used.
	RetryOn func(error) bool

	// OnRetry is called before sleeping between attempts (optional).
	OnRetry func(AttemptInfo)
}

type AttemptInfo struct {
	Attempt int
	Err     error
	Delay   time.Duration
	Elapsed time.Duration
}

func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts: 3,
		BaseDelay:   100 * time.Millisecond,
		MaxDelay:    2 * time.Second,
		Jitter:      0.2,
		RetryOn:     DefaultRetryOn,
	}
}
