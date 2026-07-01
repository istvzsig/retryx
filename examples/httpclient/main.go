package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/istvzsig/retryx"
	"github.com/istvzsig/retryx/breaker"
	"github.com/istvzsig/retryx/retry"
)

// go run ./examples/httpclient
func main() {
	httpClient := &http.Client{
		Timeout: 2 * time.Second,
	}

	policy := retry.DefaultPolicy()
	policy.MaxAttempts = 3
	policy.TimeoutPerTry = 1500 * time.Millisecond

	policy.OnRetry = func(info retry.AttemptInfo) {
		fmt.Printf("[retry] attempt=%d err=%v delay=%v\n",
			info.Attempt,
			info.Err,
			info.Delay,
		)
	}

	br := breaker.New(breaker.BreakerConfig{
		FailureThreshold: 3,
		SuccessThreshold: 2,
		OpenTimeout:      5 * time.Second,
		OnStateChange: func(from, to breaker.State) {
			fmt.Printf("[breaker] %s -> %s\n", from, to)
		},
	})

	client := retryx.Wrapper{
		Retry:   policy,
		Breaker: br,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := client.Do(ctx, func(ctx context.Context) error {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			"https://example.com",
			nil,
		)
		if err != nil {
			return err
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 500 || resp.StatusCode == 429 {
			return errors.New("retryable http status")
		}

		fmt.Printf("[ok] status=%s\n", resp.Status)
		return nil
	})

	fmt.Println("final error:", err)
}
