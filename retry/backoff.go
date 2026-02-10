package retry

import (
	"math"
	"math/rand"
	"time"
)

func backoffDelay(base, max time.Duration, attempt int, jitter float64, rnd *rand.Rand) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	if base <= 0 {
		base = 1 * time.Millisecond
	}
	if max <= 0 {
		max = base
	}
	if max < base {
		max = base
	}

	mult := math.Pow(2, float64(attempt-1))
	delay := time.Duration(float64(base) * mult)
	if delay > max {
		delay = max
	}

	if jitter <= 0 {
		return delay
	}
	if jitter > 1 {
		jitter = 1
	}

	// factor in [1-jitter, 1+jitter]
	factor := (1 - jitter) + (2 * jitter * rnd.Float64())
	j := time.Duration(float64(delay) * factor)
	if j < 0 {
		return 0
	}
	return j
}
