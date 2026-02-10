package retry

import "fmt"

type HTTPStatusError struct {
	StatusCode int
}

func (e HTTPStatusError) Error() string {
	return fmt.Sprintf("retry: http status %d", e.StatusCode)
}
