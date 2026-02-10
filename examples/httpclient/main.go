package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/istvzsig/retryx"
	"github.com/istvzsig/retryx/breaker"
	"github.com/istvzsig/retryx/retry"
)

// go run ./examples/httpclient
func main() {
	client := &http.Client{Timeout: 2 * time.Second}

	br := breaker.New(breaker.BreakerConfig{
		FailureThreshold: 3,
		SuccessThreshold: 2,
		OpenTimeout:      5 * time.Second,
		OnStateChange: func(from, to breaker.State) {
			fmt.Printf("breaker: %v -> %v\n", from, to)
		},
	})

	p := retry.DefaultPolicy()
	p.MaxAttempts = 3
	p.TimeoutPerTry = 1500 * time.Millisecond
	p.OnRetry = func(info retry.AttemptInfo) {
		fmt.Printf("retry attempt=%d err=%v delay=%v\n", info.Attempt, info.Err, info.Delay)
	}

	w := retryx.Wrapper{Retry: p, Breaker: br}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := w.Do(ctx, func(ctx context.Context) error {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", nil)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if retry.RetryHTTPStatus(resp.StatusCode) {
			return retry.HTTPStatusError{StatusCode: resp.StatusCode}
		}

		fmt.Println("ok:", resp.Status)
		return nil
	})

	fmt.Println("final err:", err)
}
