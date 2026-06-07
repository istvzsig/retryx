package main

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/istvzsig/retryx"
	"github.com/istvzsig/retryx/breaker"
	"github.com/istvzsig/retryx/retry"
)

var callCount int32

func main() {
	policy := retry.DefaultPolicy()
	policy.MaxAttempts = 3
	policy.TimeoutPerTry = 500 * time.Millisecond

	policy.OnRetry = func(info retry.AttemptInfo) {
		fmt.Printf("[retry] attempt=%d err=%v delay=%v\n",
			info.Attempt,
			info.Err,
			info.Delay,
		)
	}

	br := breaker.New(breaker.BreakerConfig{
		FailureThreshold: 2, // low threshold so we trigger quickly
		SuccessThreshold: 1,
		OpenTimeout:      3 * time.Second,
		OnStateChange: func(from, to breaker.State) {
			fmt.Printf("[breaker] %s -> %s\n", from, to)
		},
	})

	client := retryx.Wrapper{
		Retry:   policy,
		Breaker: br,
	}

	ctx := context.Background()

	// simulate unstable service
	for i := 0; i < 6; i++ {
		err := client.Do(ctx, unstableCall)

		fmt.Println("result:", err)
		fmt.Println("-----")

		time.Sleep(300 * time.Millisecond)
	}
}

func unstableCall(ctx context.Context) error {
	n := atomic.AddInt32(&callCount, 1)

	// First 4 calls always fail
	if n < 5 {
		return errors.New("simulated upstream failure")
	}

	// After that: success (if breaker allows execution)
	fmt.Println("SUCCESS: upstream call worked")
	return nil
}
