package breaker

import (
	"context"
	"sync"
	"time"
)

type BreakerConfig struct {
	Name string

	FailureThreshold int
	SuccessThreshold int
	OpenTimeout      time.Duration

	OnStateChange func(from, to State)
}

type Breaker struct {
	cfg BreakerConfig

	mu sync.Mutex

	state State

	consecutiveFailures int
	consecutiveSuccess  int

	openedAt time.Time

	// In half-open, allow only one inflight probe.
	halfOpenInFlight bool
}

func New(cfg BreakerConfig) *Breaker {
	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = 5
	}
	if cfg.SuccessThreshold <= 0 {
		cfg.SuccessThreshold = 2
	}
	if cfg.OpenTimeout <= 0 {
		cfg.OpenTimeout = 10 * time.Second
	}
	return &Breaker{cfg: cfg, state: Closed}
}

func (b *Breaker) State() State {
	var sc *stateChange

	b.mu.Lock()
	sc = b.maybeTransitionLocked(time.Now())
	st := b.state
	b.mu.Unlock()

	if sc != nil {
		sc.cb(sc.from, sc.to)
	}
	return st
}

func (b *Breaker) Do(ctx context.Context, fn func(context.Context) error) error {
	if ctx == nil {
		return context.Canceled
	}

	now := time.Now()

	if err := b.beforeCall(now); err != nil {
		return err
	}

	err := fn(ctx)

	b.afterCall(now, err)
	return err
}

func (b *Breaker) beforeCall(now time.Time) error {
	var sc *stateChange

	b.mu.Lock()
	sc = b.maybeTransitionLocked(now)

	var err error
	switch b.state {
	case Open:
		err = ErrOpen
	case HalfOpen:
		if b.halfOpenInFlight {
			err = ErrOpen
		} else {
			b.halfOpenInFlight = true
		}
	default:
		// Closed: ok
	}

	b.mu.Unlock()

	if sc != nil {
		sc.cb(sc.from, sc.to)
	}

	return err
}

func (b *Breaker) afterCall(now time.Time, err error) {
	var sc *stateChange

	b.mu.Lock()

	if b.state == HalfOpen {
		b.halfOpenInFlight = false
	}

	if err == nil {
		sc = b.onSuccessLocked(now)
	} else {
		sc = b.onFailureLocked(now)
	}

	b.mu.Unlock()

	if sc != nil {
		sc.cb(sc.from, sc.to)
	}
}

func (b *Breaker) onSuccessLocked(now time.Time) *stateChange {
	switch b.state {
	case Closed:
		b.consecutiveFailures = 0
	case HalfOpen:
		b.consecutiveSuccess++
		if b.consecutiveSuccess >= b.cfg.SuccessThreshold {
			b.consecutiveFailures = 0
			b.consecutiveSuccess = 0
			return b.recordTransitionLocked(Closed)
		}
	}
	return nil
}

func (b *Breaker) onFailureLocked(now time.Time) *stateChange {
	switch b.state {
	case Closed:
		b.consecutiveFailures++
		if b.consecutiveFailures >= b.cfg.FailureThreshold {
			b.openedAt = now
			b.consecutiveSuccess = 0
			return b.recordTransitionLocked(Open)
		}
	case HalfOpen:
		b.openedAt = now
		b.consecutiveFailures = 0
		b.consecutiveSuccess = 0
		return b.recordTransitionLocked(Open)
	}
	return nil
}

func (b *Breaker) maybeTransitionLocked(now time.Time) *stateChange {
	if b.state != Open {
		return nil
	}
	if b.openedAt.IsZero() {
		b.openedAt = now
		return nil
	}
	if now.Sub(b.openedAt) >= b.cfg.OpenTimeout {
		b.consecutiveSuccess = 0
		b.consecutiveFailures = 0
		b.halfOpenInFlight = false
		return b.recordTransitionLocked(HalfOpen)
	}
	return nil
}

func (b *Breaker) transitionLocked(to State) (from State, changed bool) {
	from = b.state
	if from == to {
		return from, false
	}
	b.state = to
	return from, true
}

// Helper function
type stateChange struct {
	cb   func(from, to State)
	from State
	to   State
}

func (b *Breaker) recordTransitionLocked(to State) (sc *stateChange) {
	from, changed := b.transitionLocked(to)
	if changed && b.cfg.OnStateChange != nil {
		return &stateChange{cb: b.cfg.OnStateChange, from: from, to: to}
	}
	return nil
}
